package router

import (
	"chef/api/page"
	"chef/core/web"
	"chef/data"
	"embed"
	"log"
	"net/http"
)

type siteService struct {
	l *log.Logger
	q *data.Queries
}

type apiService struct {
	l *log.Logger
	q *data.Queries
}

func Init(
	app *web.App,
	l *log.Logger,
	q *data.Queries,
	assets embed.FS,
) {
	site := siteService{
		l: l,
		q: q,
	}

	api := apiService{
		l: l,
		q: q,
	}

	fsHandler := http.FileServer(http.FS(assets))
	app.Mux.Handle("GET /assets/{path}", fsHandler)

	// ================================================================
	// Website Routes

	app.Handle(http.MethodGet, "/", site.HomeLoader)

	// ================================================================
	// API Routes

	app.Handle(http.MethodGet, "/api/account/init", api.AccountInit)

}

func (s siteService) HomeLoader(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	return page.Home(page.HomeData{
		Req: r,
	}).Render(ctx, w)
}
