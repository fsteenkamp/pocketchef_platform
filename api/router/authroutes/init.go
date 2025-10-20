package authroutes

import (
	"chef/core/web"
	"context"
	"net/http"
)

func (s *Service) Init(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
) error {
	acc, _, err := s.auth(ctx, web.Now(ctx), w, r)
	if err != nil {
		return err
	}

	init, err := s.Q.AccountInit(ctx, acc.ID)
	if err != nil {
		return err
	}

	return web.JSON(w, http.StatusOK, init)
}
