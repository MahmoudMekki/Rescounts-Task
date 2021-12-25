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
