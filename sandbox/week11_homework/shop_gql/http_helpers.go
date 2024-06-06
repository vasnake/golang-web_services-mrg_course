package main

import (
	"encoding/json"
	"net/http"
)

// http helpers

type AppResponseEnvelope struct {
	Body  interface{} `json:"body,omitempty"`
	Error string      `json:"error,omitempty"`
}

func writeJsonResponse(w http.ResponseWriter, body interface{}, errText string) {
	w.Header().Add("Content-Type", "application/json")
	jsonBytes, err := json.Marshal(&AppResponseEnvelope{
		Body:  body,
		Error: errText,
	})
	panicOnError("writeJsonResponse failed, can't marshal json", err)
	w.Write(jsonBytes)
	show("writeJsonResponse, written %s", string(jsonBytes))
}
