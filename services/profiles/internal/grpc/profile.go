package grpc

import (
	"context"
	"marketplace/services/profiles/internal/models"
	"marketplace/services/profiles/internal/services"
	pb "marketplace/shared/protobuf"
)

type profileServer struct {
	pb.UnimplementedProfileServiceServer
	ProfileService services.ProfileService
}

func NewProfileServer(profileService services.ProfileService) pb.ProfileServiceServer {
	return &profileServer{ProfileService: profileService}
}

func (s *profileServer) GetProfileByID(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	profile, err := s.ProfileService.GetProfileByID(req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetProfileResponse{
		Profile: &pb.Profile{
			Id:        profile.ID.String(),
			FirstName: profile.FirstName,
			LastName:  profile.LastName,
			CreatedAt: profile.CreatedAt.String(),
			UpdatedAt: profile.UpdatedAt.String(),
		},
	}, nil
}

func (s *profileServer) CreateProfile(ctx context.Context, req *pb.CreateProfileRequest) (*pb.CreateProfileResponse, error) {
	profile := &models.Profile{
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	createdProfile, err := s.ProfileService.CreateProfile(profile)
	if err != nil {
		return nil, err
	}

	return &pb.CreateProfileResponse{
		Profile: &pb.Profile{
			Id:        createdProfile.ID.String(),
			FirstName: createdProfile.FirstName,
			LastName:  createdProfile.LastName,
			CreatedAt: createdProfile.CreatedAt.String(),
			UpdatedAt: createdProfile.UpdatedAt.String(),
		},
	}, nil
}
