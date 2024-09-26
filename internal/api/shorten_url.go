package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/url"
	"shorten/internal/store"
)

type shortenURLRequest struct {
	URL string `json:"url"`
}

type shortenURLResponse struct {
	Code string `json:"code"`
}

func handleShortenURL(store store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body PostBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			sendJson(
				w,
				apiResponse{Error: "Invalid body"},
				http.StatusUnprocessableEntity,
			)
			return
		}

		if _, err := url.Parse(body.URL); err != nil {
			sendJson(
				w,
				apiResponse{Error: "Invalid url passed"},
				http.StatusBadRequest,
			)
		}

		code, err := store.SaveShortenedURL(r.Context(), body.URL)
		if err != nil {
			slog.Error("Failed to create code", "error", err)
			sendJson(
				w,
				apiResponse{Error: "Something went wrong."},
				http.StatusInternalServerError,
			)
			return
		}
		sendJson(
			w,
			apiResponse{Data: code},
			http.StatusCreated,
		)
	}
}
