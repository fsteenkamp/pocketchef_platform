package authroutes

import (
	"chef/core/randx"
	"context"
	"net/http"
)

// Possible errors
// -> invalid-email
func (s *Service) SignupWithProvider(
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
