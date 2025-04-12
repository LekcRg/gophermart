package httputils

import (
	"encoding/json"
	"net/http"
)

type MessageJSON struct {
	Message string `json:"message"`
}

type ErrorJSON struct {
	Error string `json:"error"`
}

func IsJSON(r *http.Request) bool {
	return r.Header.Get("Content-Type") == "application/json"
}

func ErrInternal(w http.ResponseWriter) {
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

func ErrMapJSON(w http.ResponseWriter, errMap map[string]string, code int) {
	w.Header().Set("Content-Type", "application/json")
	msg, err := json.Marshal(errMap)
	if err != nil {
		ErrInternalJSON(w)
	}

	w.WriteHeader(code)
	w.Write(msg)
}

func ErrJSON(w http.ResponseWriter, errMsg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	errStruct := ErrorJSON{
		Error: errMsg,
	}
	msg, err := json.Marshal(errStruct)
	if err != nil {
		ErrInternalJSON(w)
	}
	w.WriteHeader(code)
	w.Write(msg)
}

func ErrInternalJSON(w http.ResponseWriter) {
	ErrJSON(w, "Inernal server error", http.StatusInternalServerError)
}

func SuccessJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	msgStruct := MessageJSON{
		Message: "success",
	}
	msg, err := json.Marshal(msgStruct)
	if err != nil {
		ErrInternalJSON(w)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(msg)
}
