package router

import (
	"chef/api/router/adminroutes"
	"chef/api/router/authroutes"
	"chef/api/router/chefroutes"
	"chef/core/enc"
	"chef/core/web"
	"chef/data"
	"embed"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

func Init(
	app *web.App,
	l *log.Logger,
	q *data.Queries,
	assets embed.FS,
	hasher *enc.Hasher,
	googleProvider *oauth2.Config,
) {
	adminR := adminroutes.Service{
		L:      l,
		Q:      q,
		Hasher: hasher,
	}

	// publicR := publicroutes.Service{
	// 	L: l,
	// 	Q: q,
	// }

	chefRouter := chefroutes.Service{
		L:      l,
		Q:      q,
		Hasher: hasher,
	}

	authR := authroutes.Service{
		L:              l,
		Q:              q,
		Hasher:         hasher,
		GoogleProvider: googleProvider,
	}

	// ================================================================
	// SPA Routes

	fsHandler := http.FileServer(http.FS(assets))
	app.Mux.Handle("GET /assets/{path}", fsHandler)

	// ================================================================
	// Auth Routes

	app.Handle(http.MethodGet, "/api/auth/init", authR.Init)
	app.Handle(http.MethodPost, "/api/auth/signup/credentials", authR.SignupWithCredentials)
	app.Handle(http.MethodGet, "/api/auth/signup/verify", authR.SignupVerify)
	app.Handle(http.MethodPost, "/api/auth/signin/credentials", authR.SigninWithCredentials)
	app.Handle(http.MethodPost, "/api/auth/provider", authR.Provider)
	app.Handle(http.MethodGet, "/api/auth/callback/{provider}", authR.ProviderCallback)
	app.Handle(http.MethodPost, "/api/auth/signout", authR.Signout)

	// ================================================================
	// Public Routes

	// app.Handle(http.MethodGet, "/api/public/recipe/list/all", publicR.RecipeListAll)

	// ================================================================
	// Chef Routes

	app.Handle(http.MethodPost, "/api/chef/create", chefRouter.ChefProfileCreate)

	// ================================================================
	// Admin Routes

	app.Handle(http.MethodGet, "/api/admin/account/list", adminR.AccountList)
	app.Handle(http.MethodPost, "/api/admin/account/toggle/admin", adminR.AccountToggleAdmin)

}
