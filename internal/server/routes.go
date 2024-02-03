package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", s.HelloWorldHandler)
	r.Post("/log", s.SendLogHandler)
	r.Get("/logs", s.GetLogsHandler)

	return r
}

type Child struct {
	Name string `json:"name"`
}

type Parent struct {
	Childlren []Child `json:"children"`
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) SendLogHandler(w http.ResponseWriter, r *http.Request) {
	channel := "logs1"

	_, err := s.r.HSet(r.Context(), channel, "child", "child1").Result()
	if err != nil {
		log.Fatalf("error setting child in redis. Err: %v", err)
		return
	}

	resp := make(map[string]string)
	resp["message"] = "Log sent"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
		return
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) GetLogsHandler(w http.ResponseWriter, r *http.Request) {
	channel := "logs1"
	logs, err := s.r.HGetAll(r.Context(), channel).Result()
	if err != nil {
		log.Fatalf("error getting logs from redis. Err: %v", err)
		return
	}

	jsonResp, err := json.Marshal(logs)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}
