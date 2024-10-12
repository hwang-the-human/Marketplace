package main

import (
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"marketplace/services/auth/internal/config"
	"marketplace/services/auth/internal/grpc_clients"
	sharedConfig "marketplace/shared/config"
	"net/http"
	"os"
)

func main() {
	sharedConfig.InitLogrus()

	if err := godotenv.Load(".env"); err != nil {
		logrus.Fatalf("Error loading .env file")
	}

	var (
		port = os.Getenv("AUTH_PORT")
	)

	profileClient, err := grpc_clients.NewProfileClient()

	if err != nil {
		logrus.Fatalf("Error creating profile client: %v", err)
	}

	r := chi.NewRouter()

	config.InitSupertokens(r, profileClient)

	logrus.Infof("" + port)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		logrus.Fatalf("Error starting server: %v", err)
	}
}
