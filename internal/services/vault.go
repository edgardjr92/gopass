package services

import (
	"context"
	"log"
	"strings"

	"github.com/edgardjr92/gopass/internal/errors"
	"github.com/edgardjr92/gopass/internal/keys"
	"github.com/edgardjr92/gopass/internal/models"
	"github.com/edgardjr92/gopass/internal/repositories"
)

type IVaultService interface {
	// Create creates a new vault.
	// It returns the ID of the newly created vault.
	Create(ctx context.Context, name string) (uint, error)
}

type vaultService struct {
	repository repositories.IVaultRepository
}

func NewVaultService(repository repositories.IVaultRepository) *vaultService {
	return &vaultService{repository}
}

func (v *vaultService) Create(ctx context.Context, name string) (uint, error) {
	userID, ok := ctx.Value(keys.UserIDKey).(uint)

	if !ok {
		return 0, errors.UnauthorizedError("User is not authenticated")
	}

	if strings.TrimSpace(name) == "" {
		return 0, errors.BadRequestError("name is required")
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

	if err := v.repository.Save(&newVault); err != nil {
		log.Printf("Error while trying to save vault: %v", err.Error())
		return 0, err
	}

	return newVault.ID, nil
}
