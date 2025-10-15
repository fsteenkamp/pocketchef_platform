package router

import (
	"chef/core/web"
	"net/http"
)

func (s apiService) AccountInit(w http.ResponseWriter, r *http.Request) error {
	return web.JSON(w, http.StatusOK, map[string]any{
		"status": "ok",
	})
}
