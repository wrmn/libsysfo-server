package server

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Data        interface{} `json:"data,omitempty"`
	Status      int         `json:"status"`
	Reason      string      `json:"reason"`
	Description string      `json:"description"`
}

func (data response) responseFormatter(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.WriteHeader(data.Status)
	jsonResp, _ := json.Marshal(data)
	w.Write(jsonResp)
}

func unauthorizedRequest(w http.ResponseWriter, err error) {
	response{
		Status:      http.StatusUnauthorized,
		Reason:      "Unauthorized",
		Description: err.Error(),
	}.responseFormatter(w)
}

func badRequest(w http.ResponseWriter, msg string) {
	response{
		Status:      http.StatusBadRequest,
		Reason:      "Bad Request",
		Description: msg,
	}.responseFormatter(w)
}
