package authroutes

import (
	"chef/core/auth"
	"chef/core/randx"
	"chef/core/web"
	"chef/data"
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/carlmjohnson/requests"
	"github.com/jackc/pgx/v5/pgtype"
)

// Possible errors
// -> invalid-email
func (s *Service) Provider(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
) error {
	provider := r.FormValue("provider")
	state := randx.UID()

	switch provider {
	case "google":
		redirect := s.GoogleProvider.AuthCodeURL(state)
		http.Redirect(w, r, redirect, http.StatusSeeOther)
	default:
		http.Redirect(w, r, "/auth?err=invalid-provider", http.StatusSeeOther)
	}

	return nil
}

func (s *Service) ProviderCallback(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
) error {
	now := web.Now(ctx)
	provider := r.PathValue("provider")

	switch provider {
	case "google":
		return s.callbackGoogle(ctx, now, w, r)
	default:
		// TODO: add error page
	}

	return nil
}

func (s Service) callbackGoogle(ctx context.Context, now time.Time, w http.ResponseWriter, r *http.Request) error {
	s.L.Println("triggered google callback")

	authorizationCode := r.FormValue("code")

	token, err := s.GoogleProvider.Exchange(ctx, authorizationCode)
	if err != nil {
		return err
	}

	userInfo := struct {
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		FamilyName    string `json:"family_name"`
		GivenName     string `json:"given_name"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
		Sub           string `json:"sub"`
		Hd            string `json:"hd"`
	}{}

	if err := requests.
		URL("https://www.googleapis.com/oauth2/v3/userinfo").
		Method(http.MethodGet).
		Bearer(token.AccessToken).
		ToJSON(&userInfo).
		Fetch(ctx); err != nil {
		return err
	}

	var accountID string

	account, err := s.Q.AccountGetByEmail(ctx, userInfo.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// If the user does not exist, we can simply create then
			// and then carry om

			s.L.Printf("account not found, created account for %q", userInfo.Email)
			accountID = randx.UID()

			if err := s.Q.AccountCreate(ctx, data.AccountCreateParams{
				ID:       accountID,
				Email:    userInfo.Email,
				Verified: true, // SSO user is immediately verified
				Provider: pgtype.Text{
					String: "google",
					Valid:  true,
				},
				ProviderToken: pgtype.Text{
					String: token.AccessToken,
					Valid:  true,
				},
				ProviderRefreshToken: pgtype.Text{
					String: token.RefreshToken,
					Valid:  true,
				},
				ProviderLastRefresh: pgtype.Timestamp{
					Time:  now,
					Valid: true,
				},
				Picture: pgtype.Text{
					String: userInfo.Picture,
					Valid:  true,
				},
			}); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		accountID = account.ID
		// If the user does exist, we must ensure that they are an SSO user,
		// otherwise their account must be linked

		if account.Provider.Valid && account.Provider.String == "google" {
			if err := s.Q.AccountRefreshProviderDetails(ctx, data.AccountRefreshProviderDetailsParams{
				ID:        accountID,
				FirstName: pgtype.Text{String: userInfo.GivenName, Valid: true},
				LastName:  pgtype.Text{String: userInfo.FamilyName, Valid: true},
				Picture:   pgtype.Text{String: userInfo.Picture, Valid: true},
			}); err != nil {
				return err
			}
		} else {

			s.L.Println("this account is not configured for google sso, account must be linked")
			http.Redirect(w, r, "/signin?err=acc-link-required", http.StatusSeeOther)
			return nil
		}
	}

	// refresh provider information

	sessionID := randx.UID()

	expAt := now.Add(auth.SessionDuration)
	randToken := randx.UID()

	tokenHash := s.Hasher.Hash(randToken)

	if err := s.Q.SessionCreate(ctx, data.SessionCreateParams{
		ID:        sessionID,
		AccountID: accountID,
		ExpiresAt: pgtype.Timestamp{Time: expAt, Valid: true},
		TokenHash: tokenHash,
		CreatedAt: pgtype.Timestamp{Time: now, Valid: true},
	}); err != nil {
		return err
	}

	web.SetCookie(w, auth.COOKIE, randToken, now.Add(auth.SessionDuration))
	http.Redirect(w, r, "", http.StatusSeeOther)

	return nil
}
