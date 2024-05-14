package profile

import (
	"context"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
)


type service struct{
	Repo repository.ProfileStorer
}

type Service interface {
	CreateProfile(profileDetail dto.CreateProfileRequest, ctx context.Context) error
}

func NewServices(repo repository.ProfileStorer) Service{
	return &service{
		Repo: repo,
	}
}

func (profileSvc *service) CreateProfile(profileDetail dto.CreateProfileRequest, ctx context.Context) error {
	
	err := profileSvc.Repo.CreateProfile(profileDetail,ctx)
	if err != nil {
		return err
	}

	return nil
}
