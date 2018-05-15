package localhttp

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
}

const fallbackResponse = `{"status": false, "data": "An internal error occurred."}`

func respondOK(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(response{Success: true, Data: data})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fallbackResponse))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func respond(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(response{Success: false, Data: data})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fallbackResponse))
		return
	}
	w.WriteHeader(code)
	w.Write(b)
}
