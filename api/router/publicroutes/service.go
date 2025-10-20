package publicroutes

import (
	"chef/data"
	"log"
)

type Service struct {
	L *log.Logger
	Q *data.Queries
}
