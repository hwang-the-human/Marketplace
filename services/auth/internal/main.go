package main

import (
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"marketplace/services/auth/internal/config"
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

	r := chi.NewRouter()

	config.Init(r)

	logrus.Infof("" + port)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		logrus.Fatalf("Error starting server: %v", err)
	}
}
