package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type ApiConfig struct {
	ActionsSecret string
}

func main() {
	godotenv.Load()
	apiConfig := ApiConfig{}
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", apiConfig.HandleWebhook)

	actionsSecret := os.Getenv("GITHUB_ACTIONS_SECRET")
	if actionsSecret == "" {
		log.Fatal("GITHUB_ACTIONS_SECRET required")
	}
	apiConfig.ActionsSecret = actionsSecret

	newServer := http.Server{
		Handler: serveMux,
		Addr:    ":8081",
	}

	newServer.ListenAndServe()
}
