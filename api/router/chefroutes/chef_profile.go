package chefroutes

import (
	"chef/core/randx"
	"chef/core/web"
	"chef/data"
	"context"
	"net/http"
)

func (s *Service) ChefProfileCreate(
	ctx context.Context,
	w http.ResponseWriter,
	r *http.Request,
) error {
	acc, _, err := s.auth(ctx, web.Now(ctx), w, r)
	if err != nil {
		return err
	}

	body := struct {
		DisplayName string `json:"display_name"`
		Description string `json:"description"`
	}{}

	if err := web.DecodeJsonBody(w, r, s.L, &body); err != nil {
		return web.JSON(w, http.StatusBadRequest, web.JsonErr{
			Err:     web.ERR_INVALID,
			Context: "Failed to decode request.",
		})
	}

	if err := s.Q.ChefProfileCreate(ctx, data.ChefProfileCreateParams{
		ID:          randx.UID(),
		AccountID:   acc.ID,
		DisplayName: body.DisplayName,
		Description: body.Description,
	}); err != nil {
		// TODO: handle possible errors
		return err
	}

	return web.JsonOK(w)
}
