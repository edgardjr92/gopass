package services

import (
	"context"
	"log"

	"github.com/edgardjr92/gopass/internal/errors"
	"github.com/edgardjr92/gopass/internal/models"
	"github.com/edgardjr92/gopass/internal/repositories"
)

type IVaultService interface {
	// Create creates a new vault.
	// It returns the ID of the newly created vault.
	Create(ctx context.Context, name string) (int, error)
}

type vaultService struct {
	repository repositories.IVaultRepository
}

func NewVaultService(repository repositories.IVaultRepository) *vaultService {
	return &vaultService{repository}
}

func (v *vaultService) Create(ctx context.Context, name string) (int, error) {
	userID, ok := ctx.Value("user_id").(uint)

	if !ok {
		return 0, errors.UnauthorizedError("User is not authenticated")
	}

	vault, err := v.repository.FindByNameAndUserID(name, userID)

	if err != nil {
		log.Printf("Error while trying to find a vault by name,userId: %v", err.Error())
		return 0, err
	}

	if vault.ID != 0 {
		return 0, errors.ConflictError("Vault already exists")
	}

	newVault := models.Vault{
		Name:   name,
		UserID: userID,
	}

	return v.repository.Save(&newVault)
}
