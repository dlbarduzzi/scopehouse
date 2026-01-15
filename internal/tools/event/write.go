package event

import (
	"encoding/json"
	"net/http"
	"strings"
)

func WriteJson(w http.ResponseWriter, data any, status int) error {
	res, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	res = append(res, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err := w.Write(res); err != nil {
		return err
	}

	return nil
}

func WriteText(w http.ResponseWriter, status int, message string) error {
	msg := strings.TrimSpace(message)
	if msg == "" {
		msg = http.StatusText(status)
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)

	if _, err := w.Write([]byte(msg)); err != nil {
		return err
	}

	return nil
}

func WriteStatus(w http.ResponseWriter, status int) error {
	return WriteText(w, status, "")
}
