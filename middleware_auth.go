package main

import (
	"fmt"
	"github.com/dastardlyjockey/rssagg/internal/auth"
	"github.com/dastardlyjockey/rssagg/internal/database"
	"net/http"
)

type authHeader func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authHeader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Error getting API key: %s", err))
			return
		}

		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Error getting user %s", err))
			return
		}

		handler(w, r, user)
	}
}
