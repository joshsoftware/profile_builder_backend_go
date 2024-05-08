package repository

import (
	"github.com/jackc/pgx/v5"
)

type ProfileStore struct{
	db *pgx.Conn
}

type ProfileStorer interface {
	CreateProfile() error
}

func NewProfileRepo(db *pgx.Conn) ProfileStorer{
	return &ProfileStore{
        db: db,
    }
}

func (p *ProfileStore) CreateProfile() error{
	return nil
}