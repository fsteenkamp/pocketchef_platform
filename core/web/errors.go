package web

import (
	"errors"
	"net/http"
)

const (
	ERR_NOT_FOUND    = "ERR_NOT_FOUND"
	ERR_EXPIRED      = "ERR_EXPIRED"
	ERR_INVALID      = "ERR_INVALID"
	ERR_UNAUTHORISED = "ERR_UNAUTHORISED"
)

// Error is meant to be used as an early return signal. If you want to build functionality with an
// early return, you can simply return a web.Error up the call chain. At the point of needing to
// handle an error, you can choose to do nothing. The web.App will ignore an error value returned if
// it is a web.Error
var Error = errors.New("web error")

// IsError checks whether an error is of type web.Error
func IsError(err error) bool {
	return errors.Is(err, Error)
}

type JsonErr struct {
	Err     string            `json:"err"`
	Fields  map[string]string `json:"fields"`
	Context string            `json:"context"`
}

func ErrNotFound(w http.ResponseWriter, context string) error {
	return JSON(w, http.StatusNotFound, JsonErr{
		Err:     ERR_NOT_FOUND,
		Context: context,
	})
}
