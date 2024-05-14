package repository

import (
	"context"
	"errors"
	"os"

	"github.com/jackc/pgx/v5"
)

func InitializeDatabase(ctx context.Context)(*pgx.Conn, error){
	db, err := pgx.Connect(ctx, os.Getenv("PSQL_INFO"));
	if err != nil {
		return nil, errors.New("error connecting to database");
	}

	err = db.Ping(ctx)
	if err != nil{
		panic(err)
	}

	return db, nil;
}