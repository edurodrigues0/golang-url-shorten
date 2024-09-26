package api

import (
	"errors"
	"log/slog"
	"net/http"
	"shorten/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

type getShortenedURLResponse struct {
	FullURL string `json:"full_url"`
}

func handleGetShorten(store store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")
		fullUrl, err := store.GetFullURL(r.Context(), code)
		if err != nil {
			if errors.Is(err, redis.Nil) {
				sendJson(w,
					apiResponse{Error: "Code not found."},
					http.StatusNotFound,
				)
			}
			slog.Error("Failed to get code", "error", err)
			sendJson(w,
				apiResponse{Error: "Something went wrong."},
				http.StatusInternalServerError,
			)
		}

		sendJson(
			w,
			apiResponse{Data: getShortenedURLResponse{FullURL: fullUrl}},
			http.StatusOK,
		)
	}
}
