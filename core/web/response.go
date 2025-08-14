package web

import (
	"encoding/json"
	"net/http"
)

// JSON converts a Go value to JSON and sends it to the client.
func JSON(
	w http.ResponseWriter,
	statusCode int,
	data any,
) error {

	// If there is nothing to marshal then set status code and return.
	if statusCode == http.StatusNoContent || data == nil {
		w.WriteHeader(statusCode)

		return nil
	}

	// Convert the response value to JSON.
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Set the content type and headers once we know marshaling has succeeded.
	w.Header().Set("Content-Type", "application/json")

	// Write the status code to the response.
	w.WriteHeader(statusCode)

	// Write response data to response body.
	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}

func JsonOK(w http.ResponseWriter) error {
	return JSON(w, http.StatusOK, map[string]any{
		"status": "ok",
	})
}

func Text(
	w http.ResponseWriter,
	statusCode int,
	data string,
) error {

	// If there is nothing to marshal then set status code and return.
	if statusCode == http.StatusNoContent || data == "" {
		w.WriteHeader(statusCode)
		return nil
	}

	// Set the content type and headers once we know marshaling has succeeded.
	w.Header().Set("Content-Type", "text/plain")

	// Write the status code to the response.
	w.WriteHeader(statusCode)

	// Write response data to response body.
	if _, err := w.Write([]byte(data)); err != nil {
		return err
	}

	return nil
}

func NoContent(
	w http.ResponseWriter,
) error {
	w.WriteHeader(http.StatusNoContent)
	return nil
}
