package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func JSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(data)
}

func InternalError(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func BadRequest(w http.ResponseWriter, error string) {
	http.Error(w, error, http.StatusBadRequest)
}

func Accepted(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusAccepted), http.StatusAccepted)
}

func NotModified(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNotModified)
	fmt.Fprintln(w, http.StatusText(http.StatusNotModified))
}
