package app

import (
	"github.com/jackc/pgx/v5"
	"github.com/joshsoftware/profile_builder_backend_go/internal/repository"
)

type Dependencies struct{
//Repo Dependencies
}

func NewServices(db *pgx.Conn) Dependencies{
	_ = repository.NewProfileRepo(db);

	return Dependencies{

    }
}