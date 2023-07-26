package main

import (
	"encoding/json"
	"fmt"
	"github.com/dastardlyjockey/rssagg/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (apiCfg *apiConfig) handlerCreateFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameter struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameter{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}

	feedFollows, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating followed feeds: %s", err))
		return
	}

	respondWithJSON(w, 201, databaseFeedFollowToFeedFollow(feedFollows))
}

func (apiCfg *apiConfig) handlerGetFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := apiCfg.DB.GetUserFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 200, fmt.Sprintf("Failed to getting the user feed %v", err))
		return
	}

	respondWithJSON(w, 200, databaseFeedFollowsToFeedFollow(feedFollows))
}

func (apiCfg *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowID, err := uuid.Parse(chi.URLParam(r, "feedFollowID"))
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Invalid feed ID: %v", err))
		return
	}

	err = apiCfg.DB.DeleteUserFeedFollow(r.Context(), database.DeleteUserFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Failed to delete user feed follow: %v", err))
		return
	}

	respondWithJSON(w, 200, struct{}{})
}
