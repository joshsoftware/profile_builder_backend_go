package profile

import "github.com/joshsoftware/profile_builder_backend_go/internal/repository"


type service struct{
	Repo repository.ProfileStorer
}

type Service interface {

}

func NewServices(repo repository.ProfileStorer) Service{
	return &service{
		Repo: repo,
	}
}
