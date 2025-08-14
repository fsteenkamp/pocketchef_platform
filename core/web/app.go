package web

import (
	"fmt"
	"log"
	"net/http"
)

type Handler func(w http.ResponseWriter, r *http.Request) error

type App struct {
	Mux      *http.ServeMux
	l        *log.Logger
	NotFound Handler
	Origin   string
}

// NewApp creates an App that handle a set of routes for the application.
func NewApp(l *log.Logger, origin string) *App {
	return &App{
		Mux:    http.NewServeMux(),
		l:      l,
		Origin: origin,
	}
}

// ServeHTTP implements the http.Handler interface.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.Mux.ServeHTTP(w, r)
}

type logResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (c *logResponseWriter) WriteHeader(statusCode int) {
	c.statusCode = statusCode
	c.ResponseWriter.WriteHeader(statusCode)
}

// Handle sets a handler function for a given HTTP method and path pair
// to the application server mux. Optionally wrapping route-specific middleware around
// the handler
func (a *App) Handle(method string, path string, handler Handler) {
	// The function to execute for each request.
	h := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				a.l.Printf("PANIC: %s", r)
			}
		}()

		defer func() {
			if err := r.Body.Close(); err != nil {
				a.l.Printf("ERROR: closing request body: %s", err)
			}
		}()

		rw := logResponseWriter{w, http.StatusOK}

		a.l.Printf("started request %s %s", r.Method, r.URL.Path)

		if err := handler(&rw, r); err != nil {

			// if the error returned is not a web.Error,
			// then we need to return a 500
			if !IsError(err) {
				a.l.Printf("ERROR: %s", err.Error())

				// TODO: replace this with a page

				if err := JSON(w, http.StatusInternalServerError, JsonErr{
					Err:     "ERR_INTERNAL",
					Context: "Something went wrong, check server logs.",
				}); err != nil {
					a.l.Printf("ERROR: writing JSON: %s", err)
				}
			}
		}

		// this has to be after the handler call, otherwise the closure wraps up the outdated status
		// code
		a.l.Printf("ended request %s %s status=%d", r.Method, r.URL.Path, rw.statusCode)
	}

	a.Mux.HandleFunc(fmt.Sprintf("%s %s", method, path), h)
}
