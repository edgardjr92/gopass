package repositories

import "github.com/edgardjr92/gopass/internal/models"

type IUserRepository interface {
	// Save saves a user in the database.
	// It returns the ID of the newly created user.
	Save(user *models.User) (int, error)
	// FindByEmail finds a user by email.
	FindByEmail(email string) (*models.User, error)
}
