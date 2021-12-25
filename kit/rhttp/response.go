package rhttp

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// RespondJSON -- makes the tracking_resp with payload as json format
func RespondJSON(w http.ResponseWriter, httpStatusCode int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(httpStatusCode)
	w.Write(data)
}

// RespondError -- makes the error tracking_resp with payload as json format
func RespondError(w http.ResponseWriter, httpStatusCode int, message string) {
	RespondJSON(w, httpStatusCode, map[string]string{"error": message})
}

// RespondMessage -- makes the message tracking_resp with payload as json format
func RespondMessage(w http.ResponseWriter, httpStatusCode int, message string) {
	RespondJSON(w, httpStatusCode, map[string]string{"message": message})
}

// RespondString -- makes the string
func RespondString(w http.ResponseWriter, httpStatusCode int, message string) {
	w.Header().Set("Content-Type", "text/plan; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len([]byte(message))))
	w.WriteHeader(httpStatusCode)
	w.Write([]byte(message))
}
