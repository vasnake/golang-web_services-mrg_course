package httputils

import (
	"encoding/json"
	"log"
	"net/http"
)

type MyResponse struct {
	Body  interface{} `json:"body,omitempty"`
	Error string      `json:"error,omitempty"`
}

func RespJSON(w http.ResponseWriter, body interface{}) {
	w.Header().Add("Content-Type", "application/json")
	respJSON, _ := json.Marshal(&MyResponse{
		Body: body,
	})
	w.Write(respJSON)
}

func RespJSONError(w http.ResponseWriter, status int, err error, resp string) {
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	respJSON, _ := json.Marshal(&MyResponse{
		Error: resp,
	})
	w.Write(respJSON)
}
