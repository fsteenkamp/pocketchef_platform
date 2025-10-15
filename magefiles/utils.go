package main

import (
	"chef/core/conf"
	"chef/core/enc"
	"chef/data"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type deps struct {
	q      *data.Queries
	pool   *pgxpool.Pool
	hasher *enc.Hasher
}

func connect() (deps, error) {
	cfg := struct {
		Env        string
		PgHost     string `conf:"required"`
		PgUser     string `conf:"required"`
		PgPassword string `conf:"required,mask"`
		PgDb       string `conf:"required"`
		EncSecret  string `conf:"required,mask"`
		HashSecret string `conf:"required,mask"`
	}{}

	if err := conf.Parse(&cfg); err != nil {
		return deps{}, err
	}

	// "postgres://username:password@localhost:5432/database_name"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", cfg.PgUser, cfg.PgPassword, cfg.PgHost, cfg.PgDb)

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return deps{}, err
	}

	if err := pool.Ping(ctx); err != nil {
		return deps{}, fmt.Errorf("failed to reach database: %s", err)
	}

	q := data.New(pool)

	hasher := enc.NewHasher(cfg.HashSecret)

	return deps{
		q:      q,
		pool:   pool,
		hasher: hasher,
	}, nil
}
