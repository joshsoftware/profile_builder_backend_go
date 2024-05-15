package app

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/joshsoftware/profile_builder_backend_go/internal/app/profile"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
)

type Dependencies struct {
	ProfileService profile.Service
}

func NewServices(ctx context.Context, db *pgx.Conn) Dependencies {
	profileRepo := repository.NewProfileRepo(db)

	profileService := profile.NewServices(profileRepo)

	return Dependencies{
		ProfileService: profileService,
	}
}
