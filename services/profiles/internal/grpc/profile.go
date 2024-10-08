package grpc

import (
	"context"
	"marketplace/services/profiles/internal/services"
	pb "marketplace/shared/protobuf"
)

type ProfileServer struct {
	pb.UnimplementedProfileServiceServer
	ProfileService *services.ProfileService
}

func (s *ProfileServer) GetProfileByID(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	profile, err := s.ProfileService.GetProfileByID(uint(req.Id))
	if err != nil {
		return nil, err
	}

	return &pb.GetProfileResponse{
		Profile: &pb.Profile{
			Id:        uint32(profile.ID),
			FirstName: profile.FirstName,
			LastName:  profile.LastName,
			CreatedAt: profile.CreatedAt.String(),
			UpdatedAt: profile.UpdatedAt.String(),
		},
	}, nil
}
