package mocks

import (
	"github.com/edgardjr92/gopass/internal/models"
	"github.com/stretchr/testify/mock"
)

// Define mock repository
type VaultRepositoryMock struct {
	mock.Mock
}

func (m VaultRepositoryMock) FindByNameAndUserID(name string, userID uint) (*models.Vault, error) {
	args := m.Called(name, userID)
	return args.Get(0).(*models.Vault), args.Error(1)
}

func (m *VaultRepositoryMock) Save(vault *models.Vault) error {
	args := m.Called(vault)
	if len(args) > 0 {
		err, _ := args[0].(error)
		return err
	}
	return nil
}
