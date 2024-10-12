package grpc_clients

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "marketplace/shared/protobuf"
	"os"
)

type profile struct {
	client pb.ProfileServiceClient
}

func NewProfileClient() pb.ProfileServiceClient {
	profilesAddress := os.Getenv("PROFILES_GRPC_HOST") + ":" + os.Getenv("PROFILES_GRPC_PORT")
	conn, err := grpc.NewClient(profilesAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("Could not connect to profiles service: %v", err)
		return nil
	}

	client := pb.NewProfileServiceClient(conn)
	return &profile{client: client}
}

func (pc *profile) GetProfileByID(ctx context.Context, req *pb.GetProfileRequest, opts ...grpc.CallOption) (*pb.GetProfileResponse, error) {
	return pc.client.GetProfileByID(ctx, req, opts...)
}

func (pc *profile) CreateProfile(ctx context.Context, req *pb.CreateProfileRequest, opts ...grpc.CallOption) (*pb.CreateProfileResponse, error) {
	return pc.client.CreateProfile(ctx, req, opts...)
}
