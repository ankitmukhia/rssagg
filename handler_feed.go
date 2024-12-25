package main

import (
	"fmt"
	"time"
	"net/http"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ankitmukhia/rssagg/internal/database"
)

func (cfgApi *state) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User)  {
		type parameter struct {
			Name string `json:"name"`
			URL string `json:"url"`
		}

		params := parameter{}
		err := json.NewDecoder(r.Body).Decode(&params)
		if err != nil {
			responseWithError(w, 400, fmt.Sprintf("url not provided: %v", err))
			return
		}

		//? db call
		feed, err := cfgApi.db.CreateFeed(r.Context(), database.CreateFeedParams{
			ID: uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name: params.Name,
			Url: params.URL,
			UserID: user.ID,
		})

		responseWithJson(w, 201, feed)
}

func (cfgApi state) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfgApi.db.GetFeed(r.Context())
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Counldn't find feeds: %v", err))
		return
	}

	responseWithJson(w, 201, feeds)	
}
