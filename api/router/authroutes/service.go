package authroutes

import (
	"chef/core/enc"
	"chef/data"
	"log"

	"golang.org/x/oauth2"
)

type Service struct {
	L              *log.Logger
	Q              *data.Queries
	Hasher         *enc.Hasher
	GoogleProvider *oauth2.Config
}
