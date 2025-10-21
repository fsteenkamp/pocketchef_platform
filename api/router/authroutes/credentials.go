package authroutes

import (
	"chef/core/auth"
	"chef/core/pg"
	"chef/core/randx"
	"chef/core/web"
	"chef/data"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/mail"

	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

// Possible errors
// -> invalid-email
func (s *Service) SignupWithCredentials(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
) error {
	body := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := web.DecodeJsonBody(w, r, s.L, &body); err != nil {
		return web.JSON(w, http.StatusBadRequest, web.JsonErr{
			Err:     web.ERR_INVALID,
			Context: "Failed to decode request.",
		})
	}

	if _, err := mail.ParseAddress(body.Email); err != nil {
		return web.JSON(w, http.StatusBadRequest, web.JsonErr{
			Err: web.ERR_INVALID,
			Fields: map[string]string{
				"email": "Invalid email",
			},
		})
	}

	if len(body.Password) < 8 {
		return web.JSON(w, http.StatusBadRequest, web.JsonErr{
			Err: web.ERR_INVALID,
			Fields: map[string]string{
				"password": "Password must be 8 characters or more.",
			},
		})
	}

	id := randx.UID()

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	verifyCode := randx.UID()
	verifyCodeHash := s.Hasher.Hash(verifyCode)

	if err := s.Q.AccountCreate(ctx, data.AccountCreateParams{
		ID:    id,
		Email: body.Email,
		VerifyCodeHash: pgtype.Text{
			String: string(verifyCodeHash),
			Valid:  true,
		},
		PasswordHash: pgtype.Text{
			String: string(passwordHash),
			Valid:  true,
		},
	}); err != nil {
		if pg.IsErrUniqueViolation(err) {
			return web.JSON(w, http.StatusBadRequest, web.JsonErr{
				Err: web.ERR_INVALID,
				Fields: map[string]string{
					"email": "This email address is already in use.",
				},
			})
		}

		return err
	}

	// TODO: send this as an email
	fmt.Println("======================================")
	fmt.Println(verifyCode)
	fmt.Println("======================================")

	return web.JsonOK(w)
}

// Possible errors
// -> bad-password
// -> acc-not-found
// -> require-provider-[google]
func (s *Service) SigninWithCredentials(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
) error {
	now := web.Now(ctx)

	email := r.FormValue("email")
	password := r.FormValue("password")

	acc, err := s.Q.AccountGetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Redirect(w, r, r.URL.Path+"?err=acc-not-found", http.StatusNotFound)
			return nil
		}

		return err
	}

	// If there is neither a provider or a password, this is a bad state
	if !acc.PasswordHash.Valid && !acc.Provider.Valid {
		return fmt.Errorf("ERROR: bad state, account %q has neither provider nor credentials", acc.ID)
	}

	if !acc.PasswordHash.Valid {
		http.Redirect(w, r, r.URL.Path+"?err=require-provider-"+acc.Provider.String, http.StatusNotFound)
		return nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(acc.PasswordHash.String), []byte(password)); err != nil {
		http.Redirect(w, r, r.URL.Path+"?err=bad-password", http.StatusNotFound)
		return nil
	}

	sessionID := randx.UID()

	expAt := now.Add(auth.SessionDuration)
	randToken := randx.UID()

	tokenHash := s.Hasher.Hash(randToken)

	if err := s.Q.SessionCreate(ctx, data.SessionCreateParams{
		ID:        sessionID,
		AccountID: acc.ID,
		ExpiresAt: pgtype.Timestamp{Time: expAt, Valid: true},
		TokenHash: tokenHash,
		CreatedAt: pgtype.Timestamp{Time: now, Valid: true},
	}); err != nil {
		return err
	}

	web.SetCookie(w, auth.COOKIE, randToken, now.Add(auth.SessionDuration))
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

	return nil
}

func (s *Service) SignupVerify(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
) error {
	token := r.FormValue("token")
	redirect := r.FormValue("redirect") // [true | false]

	tokenHash := s.Hasher.Hash(token)

	if err := s.Q.AccountSetVerified(ctx, pgtype.Text{
		String: tokenHash,
		Valid:  true,
	}); err != nil {
		return err
	}

	if redirect == "true" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return nil
	}

	return web.JsonOK(w)
}
