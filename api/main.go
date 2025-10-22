package main

import (
	"chef/api/router"
	"chef/core/auth"
	"chef/core/conf"
	"chef/core/enc"
	"chef/core/web"
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"chef/data"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	//go:embed assets
	assets embed.FS
)

func main() {
	l := log.New(os.Stdout, "", 0)
	if err := run(l); err != nil {
		l.Printf("ERROR: %s", err)
	} else {
		l.Println("shutdown complete")
	}
}

func run(l *log.Logger) error {
	l.Println("starting server")

	ctx := context.Background()

	// ==========================================
	// CONFIG

	cfg := struct {
		Port                    string `conf:"required"`
		PgHost                  string `conf:"required"`
		PgUser                  string `conf:"required"`
		PgPassword              string `conf:"required,mask"`
		PgDb                    string `conf:"required"`
		HashSecret              string `conf:"required,mask"`
		GoogleOauthClientID     string `conf:"required"`
		GoogleOauthClientSecret string `conf:"required,mask"`
		RedirectHost            string `conf:"required"`
	}{}

	if err := conf.ParseAndPrint(&cfg); err != nil {
		return err
	}

	// ==========================================
	// Services

	hasher := enc.NewHasher(cfg.HashSecret)
	authProviderGoogle := auth.InitGoogle(
		cfg.GoogleOauthClientID,
		cfg.GoogleOauthClientSecret,
		fmt.Sprintf("%s/api/auth/callback/google", cfg.RedirectHost),
	)

	// ==========================================
	// DB

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", cfg.PgUser, cfg.PgPassword, cfg.PgHost, cfg.PgDb)

	l.Println("connecting to postgres")

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return err
	}

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to reach database: %s", err)
	}

	q := data.New(pool)

	// ==========================================
	// SERVER

	app := web.NewApp(l, "")

	router.Init(
		app,
		l,
		q,
		assets,
		hasher,
		authProviderGoogle,
	)

	addr := fmt.Sprintf("%s:%s", "0.0.0.0", cfg.Port)

	server := http.Server{
		Addr:    addr,
		Handler: app,
	}

	web.SetServerDefaults(&server)

	serverErr := make(chan error)
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		serverErr <- server.ListenAndServe()
	}()
	l.Printf("server running: %s", addr)

	select {
	case err := <-serverErr:
		return fmt.Errorf("server error: %s", err)
	case sig := <-shutdown:
		l.Println("starting server shutdown with", sig)
		defer l.Println("server shutdown complete")

		ctx, cancel := context.WithTimeout(context.Background(), web.DefaultShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			if err := server.Close(); err != nil {
				return err
			}
			return fmt.Errorf("could not gracefully shutdown server: %s", err)
		}
		return nil
	}
}
