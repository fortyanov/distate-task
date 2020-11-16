package db

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type Connection struct {
	*pgxpool.Pool
}

