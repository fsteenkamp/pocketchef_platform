package web

import (
	"net/http"
	"time"
)

const (
	DefaultReadTimeout     = time.Second * 5
	DefaultWriteTimeout    = time.Second * 10
	DefaultIdleTimeout     = time.Second * 120
	DefaultShutdownTimeout = time.Second * 20

	MegaByte = 1 << 20
)

func SetServerDefaults(server *http.Server) {
	server.MaxHeaderBytes = http.DefaultMaxHeaderBytes
	server.ReadTimeout = DefaultReadTimeout
	server.WriteTimeout = DefaultWriteTimeout
	server.IdleTimeout = DefaultIdleTimeout
}
