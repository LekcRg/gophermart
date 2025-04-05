package httputils

import (
	"encoding/json"
	"net/http"
)

type MessageJSON struct {
	Message string
}

func IsJSON(r *http.Request) bool {
	return r.Header.Get("Content-Type") == "application/json"
}

func ErrInternal(w http.ResponseWriter) {
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

func getJSON(text string) ([]byte, error) {
	msg := MessageJSON{
		Message: text,
	}
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	return msgJSON, nil
}

func ErrJSON(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "application/json")
	msg, err := getJSON(error)
	if err != nil {
		ErrInternal(w)
	}
	w.WriteHeader(code)
	w.Write(msg)
}

func ErrInternalJSON(w http.ResponseWriter) {
	ErrJSON(w, "Inernal server error", http.StatusInternalServerError)
}

func SuccessJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	msg, err := getJSON("Success")
	if err != nil {
		ErrInternal(w)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(msg)
}
