package grpc_clients

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "marketplace/shared/protobuf"
	"os"
)

type ProfileClient struct {
	client pb.ProfileServiceClient
}

func NewProfileClient() *ProfileClient {
	profilesAddress := os.Getenv("PROFILES_GRPC_HOST") + ":" + os.Getenv("PROFILES_GRPC_PORT")
	conn, err := grpc.NewClient(profilesAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("Could not connect to profiles service: %v", err)
		return nil
	}

	client := pb.NewProfileServiceClient(conn)
	return &ProfileClient{client: client}
}

func (pc *ProfileClient) GetProfileByID(id string) (*pb.Profile, error) {
	req := &pb.GetProfileRequest{Id: id}
	res, err := pc.client.GetProfileByID(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return res.Profile, nil
}
