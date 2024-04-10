package net

import (
	"encoding/json"
	"net/http"
)

func JsonWrite(statusCode int, w http.ResponseWriter, r *http.Request, i interface{}) {
	// Set the status
	w.WriteHeader(statusCode)

	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Marshall the interface into JSON
	res, err := json.Marshal(i)
	if err != nil {
		// If there was an error, write the error to the response
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write the JSON to the response
	w.Write(res)
}

func JsonWriteError(w http.ResponseWriter, r *http.Request, i interface{}) {
	JsonWrite(http.StatusBadRequest, w, r, i)
}

func JsonWriteOk(w http.ResponseWriter, r *http.Request, i interface{}) {
	JsonWrite(http.StatusOK, w, r, i)
}
