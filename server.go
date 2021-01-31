package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
)

type server struct {
	router  *chi.Mux
	entries map[string]textData
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", `{"status":"OK"}`)
}

func sendText(s *server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var newText textData
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		json.Unmarshal(reqBody, &newText)

		var id string
		for {
			id = uuid.New().String()
			if _, ok := s.entries[id]; !ok {
				break
			}
		}
		s.entries[id] = newText
		w.WriteHeader(http.StatusCreated)

		Conv, _ := json.Marshal(textID{ID: id})
		fmt.Fprintf(w, "%s", string(Conv))
	}
}

func pasteText(s *server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id := chi.URLParam(r, "id")
		text, ok := s.entries[id]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)

		newEntry := textEntry{
			textID{
				id,
			},
			textData{
				Description: text.Description,
				Content:     text.Content,
			},
		}

		conv, _ := json.Marshal(newEntry)
		fmt.Fprintf(w, "%s", string(conv))
	}
}

func newServer() *server {

	var s server
	s.entries = make(map[string]textData)

	s.router = chi.NewRouter()
	s.router.Use(middleware.Logger)

	s.router.Get("/healthy", healthCheck)
	s.router.Post("/", sendText(&s))
	s.router.Get("/paste/{id}", pasteText(&s))

	return &s
}

func (s *server) Run(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
