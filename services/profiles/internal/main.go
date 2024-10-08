package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	ps "marketplace/services/profiles/internal/grpc"
	"marketplace/services/profiles/internal/repositories"
	"marketplace/services/profiles/internal/services"
	"marketplace/shared/config"
	"marketplace/shared/db"
	"marketplace/shared/kafka"
	"marketplace/shared/models"
	pb "marketplace/shared/protobuf"
	"net"
	"os"
)

func main() {
	config.InitLogrus()

	if err := godotenv.Load(".env"); err != nil {
		logrus.Fatal("Error loading .env file")
	}

	var (
		grpcPort   = os.Getenv("PROFILES_GRPC_PORT")
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

	profileRepository := &repositories.ProfileRepository{DB: database}
	profileService := &services.ProfileService{ProfileRepository: profileRepository}

	c := cron.New()
	if _, err := c.AddFunc("@every 1s", func() {
		kafka.ProcessOutboxMessages(database, kafka.Producer)
	}); err != nil {
		logrus.Fatalf("Failed to add cron job: %v", err)
	}

	c.Start()
	defer c.Stop()

	logrus.Infof("Successfully started outbox message processing every 1s")

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		logrus.Fatalf("Failed to listen on port %s: %v", grpcPort, err)
	}

	grpcServer := grpc.NewServer()
	profileGrpcServer := &ps.ProfileServer{ProfileService: profileService}

	pb.RegisterProfileServiceServer(grpcServer, profileGrpcServer)
	reflection.Register(grpcServer)

	logrus.Infof("Starting gRPC server on port %s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		logrus.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
