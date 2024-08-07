package routers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type errorJsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type Helpers struct {
}

// ReadJSON tries to read the body of a request and converts it into JSON
func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // one megabyte

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

// WriteJSON takes a response status code and arbitrary data and writes a json response to the client
func WriteJSON(w http.ResponseWriter, status int, data any, headers map[string]string) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for key, value := range headers {
		w.Header().Set(key, value)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

// ErrorJSON takes an error, and optionally a response status code, and generates and sends a json error response
func ErrorJSON(w http.ResponseWriter, status int, err error, headers map[string]string) error {
	statusCode := http.StatusBadRequest
	if status > 0 {
		statusCode = status
	}

	for key, value := range headers {
		w.Header().Set(key, value)
	}

	var payload errorJsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return WriteJSON(w, statusCode, payload, nil)
}
