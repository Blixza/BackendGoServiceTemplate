package user_domain

import (
	"database/sql"
	"time"

	uuid "github.com/google/uuid"
)

type Model struct {
	ID        uuid.UUID
	Nickname  string
	Password  string
	Discord   string
	Email     string
	Balance   int
	Towns     sql.NullString
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}
