package services

import (
	"context"
	"log"

	"github.com/edgardjr92/gopass/internal/cerrors"
	"github.com/edgardjr92/gopass/internal/keys"
	"github.com/edgardjr92/gopass/internal/models"
	"github.com/edgardjr92/gopass/internal/repositories"
	"github.com/edgardjr92/gopass/internal/utils"
)

type IVaultService interface {
	// Create creates a new vault.
	// It returns the ID of the newly created vault.
	Create(ctx context.Context, name string) (uint, error)
	// GetAll returns all vaults from a user.
	GetAll(ctx context.Context, userID uint) ([]models.VaultDetail, error)
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
		return 0, cerrors.UnauthorizedError("user is not authenticated")
	}

	if utils.IsBlank(name) {
		return 0, cerrors.BadRequestError("name is required")
	}

	vault, err := v.repository.FindByNameAndUserID(ctx, name, userID)

	if err != nil {
		log.Printf("error while trying to find a vault by name,userId: %v", err.Error())
		return 0, err
	}

	if vault.ID != 0 {
		return 0, cerrors.ConflictError("vault already exists")
	}

	newVault := models.Vault{
		Name:   name,
		UserID: userID,
	}

	if err := v.repository.Save(ctx, &newVault); err != nil {
		log.Printf("error while trying to save vault: %v", err.Error())
		return 0, err
	}

	return newVault.ID, nil
}

func (v *vaultService) GetAll(ctx context.Context) ([]models.VaultDetail, error) {
	userID, ok := ctx.Value(keys.UserIDKey).(uint)

	if !ok {
		return []models.VaultDetail{}, cerrors.UnauthorizedError("user is not authenticated")
	}

	vaults, err := v.repository.FindByUserID(ctx, userID)

	if err != nil {
		log.Printf("error while trying to find all vaults by userId: %v", err.Error())
		return []models.VaultDetail{}, err
	}

	details := utils.Map(vaults, func(v models.Vault) models.VaultDetail {
		return models.VaultDetail{
			ID:     v.ID,
			Name:   v.Name,
			UserID: v.UserID,
		}
	})

	return details, nil
}
