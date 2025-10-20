package router

import (
	"chef/api/router/adminroutes"
	"chef/api/router/authroutes"
	"chef/core/enc"
	"chef/core/web"
	"chef/data"
	"embed"
	"log"
	"net/http"
)

func Init(
	app *web.App,
	l *log.Logger,
	q *data.Queries,
	assets embed.FS,
	hasher *enc.Hasher,
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

	authR := authroutes.Service{
		L:      l,
		Q:      q,
		Hasher: hasher,
	}

	fsHandler := http.FileServer(http.FS(assets))
	app.Mux.Handle("GET /assets/{path}", fsHandler)

	// ================================================================
	// Auth Routes

	app.Handle(http.MethodGet, "/api/auth/init", authR.Init)

	// ================================================================
	// Public Routes

	// app.Handle(http.MethodGet, "/api/account/init", api.AccountInit)

	// ================================================================
	// Admin Routes

	app.Handle(http.MethodGet, "/api/admin/account/list", adminR.AccountList)

}
