package router

import (
	"chef/api/page"
	"chef/core/web"
	"chef/data"
	"embed"
	"log"
	"net/http"
)

type service struct {
	l *log.Logger
	q *data.Queries
}

func Init(
	app *web.App,
	l *log.Logger,
	q *data.Queries,
	assets embed.FS,
) {
	s := service{
		l: l,
		q: q,
	}

	fsHandler := http.FileServer(http.FS(assets))

	app.Mux.Handle("GET /assets/{path}", fsHandler)
	app.Handle(http.MethodGet, "/", s.HomeLoader)
}

func (s service) HomeLoader(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	return page.Home(page.HomeData{
		Req: r,
	}).Render(ctx, w)
}
