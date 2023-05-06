package mocks

import (
	"context"

	"github.com/edgardjr92/gopass/internal/models"
	"github.com/stretchr/testify/mock"
)

// Define mock repository
type VaultRepositoryMock struct {
	mock.Mock
}

func (m VaultRepositoryMock) FindByUserID(ctx context.Context, userID uint) ([]models.Vault, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Vault), args.Error(1)
}

func (m VaultRepositoryMock) FindByNameAndUserID(ctx context.Context, name string, userID uint) (*models.Vault, error) {
	args := m.Called(ctx, name, userID)
	return args.Get(0).(*models.Vault), args.Error(1)
}

func (m *VaultRepositoryMock) Save(ctx context.Context, vault *models.Vault) error {
	args := m.Called(ctx, vault)
	if len(args) > 0 {
		err, _ := args[0].(error)
		return err
	}
	return nil
}
