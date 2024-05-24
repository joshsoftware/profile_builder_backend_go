package app

import (
	"github.com/jackc/pgx/v5"
	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
)

// Dependencies struct holds dependencies required by the application.
type Dependencies struct {
	ProfileService service.Service
}

// NewServices initializes and returns service dependencies.
func NewServices(db *pgx.Conn) Dependencies {

	profileService := service.NewServices(db)

	return Dependencies{
		ProfileService: profileService,
	}
}
