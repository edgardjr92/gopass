package repositories

import (
	"context"

	"github.com/edgardjr92/gopass/internal/models"
)

type IUserRepository interface {
	// Save saves a user in the database.
	Save(ctx context.Context, user *models.User) error
	// FindByEmail finds a user by email.
	FindByEmail(ctx context.Context, email string) (*models.User, error)
}
