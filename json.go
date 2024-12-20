package main

import (
 	"log"	
	"net/http" 
	"encoding/json" 
)

func responseWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)	
	}

	/*
	* Error field, I want key for field to be error
	* { "error": "Something went wrong!" }
	*/
	type errorMsg struct {
		Error string `json:"error"`
	}

	responseWithJson(w, code, errorMsg{
		Error: msg,
	})
}

func responseWithJson(w http.ResponseWriter, code int, payload interface{}) {
	json, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON %v", payload)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(json)
}
