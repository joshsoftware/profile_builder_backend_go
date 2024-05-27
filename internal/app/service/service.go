package service

import (
	"github.com/jackc/pgx/v5"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
)

type service struct {
	userLoginRepo repository.UserStorer
}

type Service interface {
	GenerateLoginToken(email string) (string, error)
}

func NewServices(db *pgx.Conn) Service {
	userLoginRepo := repository.NewUserLoginRepo(db)
	return &service{
		userLoginRepo: userLoginRepo,
	}
}
