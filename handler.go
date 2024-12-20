package main

import (
	"net/http"
)	

func handlerErr(w http.ResponseWriter, r *http.Request) {
	responseWithError(w, 400, "Something went wrong")
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	responseWithJson(w, 200, struct{}{})
}
