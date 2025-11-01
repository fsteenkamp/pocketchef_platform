package chefroutes

import (
	"chef/core/enc"
	"chef/data"
	"log"
)

type Service struct {
	L      *log.Logger
	Q      *data.Queries
	Hasher *enc.Hasher
}
