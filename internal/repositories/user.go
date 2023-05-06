package repositories

import "github.com/edgardjr92/gopass/internal/models"

type IUserRepository interface {
	// Save saves a user in the database.
	Save(user *models.User) error
	// FindByEmail finds a user by email.
	FindByEmail(email string) (*models.User, error)
}
