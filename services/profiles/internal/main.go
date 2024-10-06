package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"marketplace/services/profiles/internal/handlers"
	"marketplace/services/profiles/internal/routes"
	"marketplace/services/profiles/internal/services"
	"marketplace/shared/config"
	"marketplace/shared/db"
	"marketplace/shared/kafka"
	"marketplace/shared/models"
	"net/http"
	"os"
)

func main() {
	config.InitLogrus()

	if err := godotenv.Load(".env"); err != nil {
		logrus.Fatal("Error loading .env file")
	}

	var (
		port       = os.Getenv("PROFILES_PORT")
		dbUser     = os.Getenv("DB_USER")
		dbPassword = os.Getenv("DB_PASSWORD")
		dbHost     = os.Getenv("DB_HOST")
		dbPort     = os.Getenv("DB_PORT")
		dbName     = os.Getenv("PROFILES_DB_NAME")
		kafkaHost  = os.Getenv("KAFKA_HOST")
	)

	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbPort, dbName)

	db.InitDB(dsn, &models.OutboxMessage{})
	defer db.CloseDB()

	database := db.GetDB()

	brokers := []string{kafkaHost}
	kafka.InitKafkaProducer(brokers)
	defer kafka.CloseKafkaProducer()

	profileService := &services.ProfileService{DB: database}
	profileHandler := &handlers.ProfileHandler{ProfileService: profileService}

	r := chi.NewRouter()
	r.Mount("/api", routes.ProfileRouter(profileHandler))

	c := cron.New()
	_, err := c.AddFunc("@every 1s", func() {
		kafka.ProcessOutboxMessages(database, kafka.Producer)
	})
	logrus.Infof("Successfully started outbox message processing every 1s")
	if err != nil {
		logrus.Fatalf("Failed to add cron job: %v", err)
	}

	c.Start()
	defer c.Stop()

	logrus.Infof("Starting profiles service on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		logrus.Fatalf("Error starting server: %v", err)
	}
}
