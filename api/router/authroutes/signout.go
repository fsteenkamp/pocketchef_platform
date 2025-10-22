package authroutes

import (
	"chef/core/auth"
	"chef/core/web"
	"chef/data"
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
)

func (s Service) Signout(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
) error {
	now := web.Now(ctx)

	account, session, err := s.auth(ctx, now, w, r)
	if err != nil {
		return err
	}

	if err := s.Q.SessionInvalidate(ctx, data.SessionInvalidateParams{
		InvalidatedAt: pgtype.Timestamp{
			Time:  now,
			Valid: true,
		},
		ID:        session.ID,
		AccountID: account.ID,
	}); err != nil {
		return err
	}

	web.DeleteCookie(w, auth.COOKIE, "/", now)
	http.Redirect(w, r, "/signin", http.StatusSeeOther)

	return nil
}
