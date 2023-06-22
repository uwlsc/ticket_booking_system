package user

import (
	"ticket-api/pkg"
)

// Repository database structure
type Repository struct {
	pkg.Database
	logger pkg.Logger
}

// NewRepository creates a new user repository
func NewRepository(db pkg.Database, logger pkg.Logger) Repository {
	return Repository{db, logger}
}
