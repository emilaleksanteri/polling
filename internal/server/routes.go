package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", s.HelloWorldHandler)
	r.Post("/log/{channel}", s.SendLogHandler)
	r.Get("/logs/{channel}", s.GetLogsHandler)

	return r
}

type Log struct {
	Domain     string    `json:"domain" redis:"domain"`
	RequestAt  time.Time `json:"request_at" redis:"request_at"`
	IsMyDomain bool      `json:"is_my_domain" redis:"is_my_domain"`
}

type Event struct {
	Id  string `json:"id" redis:"id"`
	Log Log    `json:"log" redis:"log"`
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
	channel := chi.URLParam(r, "channel")
	if channel == "" {
		w.Write([]byte("channel is required"))
		return
	}

	var input struct {
		Domain string `json:"domain"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Fatalf("error decoding JSON. Err: %v", err)
		return
	}

	myDomain := "localhost"
	event := Log{
		Domain:     input.Domain,
		RequestAt:  time.Now(),
		IsMyDomain: input.Domain == myDomain,
	}

	jsonEvent, err := json.Marshal(event)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
		return
	}

	for i := 0; i < 100; i++ {
		_, err = s.r.LPush(r.Context(), channel, jsonEvent).Result()
		if err != nil {
			log.Fatalf("error setting child in redis. Err: %v", err)
			return
		}
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
	channel := chi.URLParam(r, "channel")
	if channel == "" {
		w.Write([]byte("channel is required"))
		return
	}

	logs, err := s.r.LRange(r.Context(), channel, 0, -1).Result()
	if err != nil {
		log.Fatalf("error getting logs from redis. Err: %v", err)
		return
	}

	var logsS []Log
	for _, log := range logs {
		var logS Log
		err := json.Unmarshal([]byte(log), &logS)
		if err != nil {
			fmt.Printf("error handling JSON unmarshal. Err: %v", err)
			return
		}

		logsS = append(logsS, logS)
	}

	resp := make(map[string]any)
	resp["len"] = len(logsS)
	resp["logs"] = logsS

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}
