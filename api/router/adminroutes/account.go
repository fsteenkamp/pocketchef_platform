package adminroutes

import (
	"chef/core/web"
	"chef/data"
	"context"
	"database/sql"
	"errors"
	"net/http"
)

func (s Service) AccountList(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
) error {
	now := web.Now(ctx)
	if _, _, err := s.auth(ctx, now, w, r); err != nil {
		return err
	}

	data, err := s.Q.AccountList(ctx)
	if err != nil {
		return err
	}

	return web.JSON(w, http.StatusOK, data)
}

func (s Service) AccountToggleAdmin(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	now := web.Now(ctx)

	_, _, err := s.auth(ctx, now, w, r)
	if err != nil {
		return err
	}

	body := struct {
		AccountID string `json:"account_id"`
		IsAdmin   bool   `json:"is_admin"`
	}{}

	if err := web.DecodeJsonBody(w, r, s.L, &body); err != nil {
		return err
	}

	targetAccount, err := s.Q.AccountGetByID(ctx, body.AccountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return web.ErrNotFound(w, "Target account not found.")
		}
		return err
	}

	// cannot disable admin privelages for root account
	if targetAccount.IsRoot {
		return web.JSON(w, http.StatusBadRequest, web.JsonErr{
			Err:     web.ERR_INVALID,
			Context: "Target account is root account.",
			Fields:  map[string]string{},
		})
	}

	if err := s.Q.AccountSetAdmin(ctx, data.AccountSetAdminParams{
		ID:      body.AccountID,
		IsAdmin: body.IsAdmin,
	}); err != nil {
		return err
	}

	return web.JSON(w, http.StatusOK, map[string]any{
		"status": "ok",
	})
}
