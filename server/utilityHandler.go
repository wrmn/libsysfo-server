package server

import (
	"encoding/json"
	"net/http"
)

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

func intServerError(w http.ResponseWriter, err error) {
	response{
		Status:      http.StatusInternalServerError,
		Reason:      "Internal Server Error",
		Description: err.Error(),
	}.responseFormatter(w)
}

func handleNotFound() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		response{
			Status:      http.StatusNotFound,
			Reason:      "Not Found",
			Description: "Request not found",
		}.responseFormatter(w)
	})
}

func handleNotAllowed() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		response{
			Status:      http.StatusMethodNotAllowed,
			Reason:      "Method Not Allowed",
			Description: "Request not allowed with this method",
		}.responseFormatter(w)
	})
}
