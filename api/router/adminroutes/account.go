package adminroutes

import (
	"chef/core/web"
	"context"
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
