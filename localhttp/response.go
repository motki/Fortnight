package localhttp

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Success bool        `json:"success"`
	Status  int         `json:"status"`
	Data    interface{} `json:"data,omitempty"`
}

const fallbackResponse = `{"success": false, "status": 500, "data": "An internal error occurred."}`

func respondOK(w http.ResponseWriter, data interface{}) {
	respond(w, http.StatusOK, data)
}

func respond(w http.ResponseWriter, code int, data interface{}) {
	b, err := json.Marshal(response{
		Success: code >= 200 && code < 400,
		Status:  code,
		Data:    data})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fallbackResponse))
		return
	}
	w.WriteHeader(code)
	w.Write(b)
}
