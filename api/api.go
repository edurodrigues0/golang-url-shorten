package api

import (
	"encoding/json"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewHandler(db map[string]string) http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	r.Post("/api/shorten", handleCreateShorten(db))
	r.Get("/{code}", handleGetShorten(db))

	return r
}

type PostBody struct {
	URL string `json:"url"`
}

type Response struct {
	Error string `json:"error,omitempty"`
	Data  any    `json:"data,omitempty"`
}

func sendJson(w http.ResponseWriter, resp Response, status int) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(resp)
	if err != nil {
		slog.Error("Failed to marshal json data", "error", err)
		sendJson(
			w,
			Response{Error: "Something went wrong"},
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(status)
	if _, err := w.Write(data); err != nil {
		slog.Error("Failed to write response to client", "error", err)
		return
	}
}

func handleCreateShorten(db map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body PostBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			sendJson(
				w,
				Response{Error: "Invalid body"},
				http.StatusUnprocessableEntity,
			)
			return
		}

		if _, err := url.Parse(body.URL); err != nil {
			sendJson(
				w,
				Response{Error: "Invalid url passed"},
				http.StatusBadRequest,
			)
		}

		code := genCode()
		db[code] = body.URL
		sendJson(
			w,
			Response{Data: code},
			http.StatusCreated,
		)
	}
}

const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func genCode() string {
	const n = 8
	byts := make([]byte, n)

	for i := range n {
		byts[i] = characters[rand.IntN(len(characters))]
	}

	return string(byts)
}

func handleGetShorten(db map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")
		data, ok := db[code]
		if !ok {
			http.Error(w, "Url not found", http.StatusNotFound)
			return
		}

		http.Redirect(w, r, data, http.StatusPermanentRedirect)
	}
}
