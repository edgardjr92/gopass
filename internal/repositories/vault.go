package repositories

import "github.com/edgardjr92/gopass/internal/models"

type IVaultRepository interface {
	// Store stores a new vault.
	// It returns the ID of the newly created vault.
	Save(vault *models.Vault) (int, error)
	// Find a vault by name and user ID.
	FindByNameAndUserID(name string, userID uint) (*models.Vault, error)
}
