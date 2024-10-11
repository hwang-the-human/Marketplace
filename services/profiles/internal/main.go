package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/supertokens/supertokens-golang/recipe/jwt"
	"github.com/supertokens/supertokens-golang/supertokens"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	ps "marketplace/services/profiles/internal/grpc"
	"marketplace/services/profiles/internal/repositories"
	"marketplace/services/profiles/internal/services"
	"marketplace/shared/config"
	"marketplace/shared/db"
	"marketplace/shared/interceptors"
	"marketplace/shared/kafka"
	"marketplace/shared/models"
	"marketplace/shared/outbox"
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
		stUri      = os.Getenv("ST_URI")
		authUri    = os.Getenv("AUTH_URI")
	)

	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbPort, dbName)

	database, err := db.NewPostgresDB(dsn)
	if err != nil {
		logrus.Fatalf("Failed to initialize database: %v", err)
	}
	defer func(database db.Database) {
		err := database.CloseDB()
		if err != nil {
			logrus.Errorf("Failed to close database: %v", err)
		}
	}(database)

	err = database.Migrate(&models.OutboxMessage{})
	if err != nil {
		logrus.Fatalf("Migration failed: %v", err)
	}

	if err := supertokens.Init(supertokens.TypeInput{
		AppInfo: supertokens.AppInfo{
			APIDomain: authUri,
		},
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: stUri,
		},
		RecipeList: []supertokens.Recipe{
			jwt.Init(nil),
		},
	}); err != nil {
		logrus.Fatalf("Error initializing Supertokens: %v", err)
	}

	brokers := []string{kafkaHost}
	kafkaProducer, err := kafka.NewProducer(brokers)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer func(kafkaProducer kafka.Producer) {
		err := kafkaProducer.Close()
		if err != nil {
			log.Printf("Failed to close Kafka producer: %v", err)
		}
	}(kafkaProducer)

	outboxService := outbox.NewOutbox(database, kafkaProducer)

	c := cron.New()
	if _, err := c.AddFunc("@every 1s", func() {
		if err := outboxService.ProcessOutboxMessages(); err != nil {
			logrus.Errorf("Failed to process outbox messages: %v", err)
		}
	}); err != nil {
		logrus.Fatalf("Failed to add cron job: %v", err)
	}

	c.Start()
	defer c.Stop()

	logrus.Infof("Successfully started outbox message processing every 1s")

	profileRepository := repositories.NewProfileRepository(database)
	profileService := services.NewProfileService(profileRepository)

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		logrus.Fatalf("Failed to listen on port %s: %v", grpcPort, err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptors.JWTAuth))
	profileGrpcServer := ps.NewProfileServer(profileService)

	pb.RegisterProfileServiceServer(grpcServer, profileGrpcServer)
	reflection.Register(grpcServer)

	logrus.Infof("Starting Profile GRPC server on port %s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		logrus.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
