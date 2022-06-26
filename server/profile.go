package server

import (
	"encoding/json"
	"net/http"
)

func loginGoogle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode("well")
}
