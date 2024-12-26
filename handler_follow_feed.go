package main

import (
	"fmt"
	"time"
	"net/http"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/go-chi/chi/v5"

	"github.com/ankitmukhia/rssagg/internal/database"
)

func (cfgApi state) handlerFollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameter struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	params := parameter{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("feed id not provided: %v", err))
		return
	}

	//? database call
	feedFollow, err := cfgApi.db.CreateFollowFeed(r.Context(), database.CreateFollowFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: params.FeedID,
	})

	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("unable to create follow feed: %v", err))
		return
	}
	
	responseWithJson(w, 201, feedFollow)
}

func (cfgApi state) handlerGetFollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	//? db call
	allFollowFeed, err := cfgApi.db.GetFollowFeed(r.Context(), user.ID)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("couldn't find feeds: %v", err))
		return
	}

	responseWithJson(w, 201, allFollowFeed)
}

func (cfgApi state) handlerDeleteFollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	//? db call
	paramsId := chi.URLParam(r, "feedFollowId")
	feedFollowId, err := uuid.Parse(paramsId)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldn't parse follow Id: %v", err))
		return
	}
	
	err = cfgApi.db.DeleteFollowFeed(r.Context(), database.DeleteFollowFeedParams{
		ID: feedFollowId,
		UserID: user.ID,
	})
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldn't delete feed follow: %v", err))
		return
	}
	responseWithJson(w, 200, struct{}{})
}
