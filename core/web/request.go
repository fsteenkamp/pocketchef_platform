package web

import (
	"encoding/json"
	"log"
	"net/http"
)

func DecodeJsonBody(w http.ResponseWriter, r *http.Request, l *log.Logger, body any) error {
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		l.Printf("BAD REQUEST: %s", err)

		if err.Error() == "EOF" {
			if err := JSON(w, http.StatusBadRequest, JsonErr{
				Err:     "ERR_BAD_REQUEST",
				Context: "Missing or incomplete request JSON body.",
				Fields:  map[string]string{},
			}); err != nil {
				l.Printf("ERROR writing resp: %s", err)
			}

		} else {
			if err := JSON(w, http.StatusBadRequest, JsonErr{
				Err:     "ERR_BAD_REQUEST",
				Context: "Failed to decode request.",
				Fields:  map[string]string{},
			}); err != nil {
				l.Printf("ERROR writing resp: %s", err)
			}
		}

		return Error
	}

	return nil
}
