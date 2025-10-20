package adminroutes

import (
	"chef/core/auth"
	"chef/core/web"
	"chef/data"
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func respondInvalidToken(w http.ResponseWriter) (data.AccountGetByIDRow, data.Session, error) {
	web.JSON(w, http.StatusUnauthorized, web.JsonErr{
		Err:     "ERR_UNAUTHENTICATED",
		Context: "Invalid auth token.",
		Fields:  map[string]string{},
	})

	return data.AccountGetByIDRow{}, data.Session{}, web.Error
}

func (s *Service) auth(
	ctx context.Context,
	now time.Time,
	w http.ResponseWriter,
	r *http.Request,
) (data.AccountGetByIDRow, data.Session, error) {
	c, err := r.Cookie(auth.COOKIE)
	var token string
	if err == nil {
		token = c.Value
	} else {
		// check for bearer token next

		bearer := r.Header.Get("Authorization")
		if bearer == "" {
			web.JSON(w, http.StatusUnauthorized, web.JsonErr{
				Err:     "ERR_UNAUTHENTICATED",
				Context: "The auth token is missing from the request.",
				Fields:  map[string]string{},
			})

			return data.AccountGetByIDRow{}, data.Session{}, web.Error
		}

		token = strings.TrimPrefix(bearer, "Bearer ")
	}

	tokenHash := s.Hasher.Hash(token)

	session, err := s.Q.SessionGetFromTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return respondInvalidToken(w)
		}

		return data.AccountGetByIDRow{}, data.Session{}, err
	}

	if session.Invalidated {
		return respondInvalidToken(w)
	}

	acc, err := s.Q.AccountGetByID(ctx, session.AccountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return respondInvalidToken(w)
		}

		return data.AccountGetByIDRow{}, data.Session{}, err
	}

	if !acc.IsAdmin {
		web.JSON(w, http.StatusForbidden, web.JsonErr{
			Err:     "ERR_FORBIDDEN",
			Context: "Account does not have admin role.",
			Fields:  map[string]string{},
		})

		return data.AccountGetByIDRow{}, data.Session{}, web.Error
	}

	// -> if the account LastActive has never been set before we set it OR
	// -> if the LastActive has not been set in the last 5 minutes we set it
	if !acc.LastActive.Valid || now.After(acc.LastActive.Time.Add(time.Minute*5)) {
		go func() {
			ctx := context.Background()

			if err := s.Q.AccountSetLastActive(ctx, data.AccountSetLastActiveParams{
				LastActive: pgtype.Timestamp{
					Time:  now,
					Valid: true,
				},
				ID: acc.ID,
			}); err != nil {
				s.L.Printf("ERROR: setting account.last_active: %s", err)
			}
		}()
	}

	return acc, session, nil
}
