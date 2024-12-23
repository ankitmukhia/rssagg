package main

import (
	"net/http"
	"fmt"

	"github.com/ankitmukhia/rssagg/internal/database"
	"github.com/ankitmukhia/rssagg/internal/auth"
)

type handlers func(w http.ResponseWriter, r *http.Request, database database.User)

func (apiCfg *state) middlewareAuth(next  handlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetApiKey(r.Header)
		if err != nil {
			responseWithError(w, 403, fmt.Sprintf("unauthorized user %v", err))
			return
		}

		user, err := apiCfg.db.GetUser(r.Context(), token)
		if err != nil {
			responseWithJson(w, 400, fmt.Sprintf("Couldn't get user %v", err))
			return
		}

		next(w, r, user)	
	}
}
