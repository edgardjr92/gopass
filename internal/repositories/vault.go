package repositories

import (
	"context"

	"github.com/edgardjr92/gopass/internal/models"
)

type IVaultRepository interface {
	// Store stores a new vault.
	Save(ctx context.Context, vault *models.Vault) error
	// Find a vault by name and user ID.
	FindByNameAndUserID(ctx context.Context, name string, userID uint) (*models.Vault, error)
	// FindByUserID returns all vaults from a user.
	FindByUserID(ctx context.Context, userID uint) ([]models.Vault, error)
}
