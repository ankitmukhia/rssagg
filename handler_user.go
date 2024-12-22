package main

import (
	"encoding/json" 
	"fmt"
	"net/http" 
	"time"
	
	"github.com/google/uuid"

	"github.com/ankitmukhia/rssagg/internal/database"
	"github.com/ankitmukhia/rssagg/internal/auth"
)

func (apiCfg *state) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	params := parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Error parsing JOSN %v", err))
		return
	}

	user, err := apiCfg.db.CreateUser(r.Context(), database.CreateUserParams{
	 	ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: params.Name,
	})

	if err != nil {
		responseWithJson(w, 500, fmt.Sprintf("Couldn't create user %v", err))
		return
	}

	responseWithJson(w, 201, user)
}

func (apiCfg state) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	// auth endpoint
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

	responseWithJson(w, 201, user)
}