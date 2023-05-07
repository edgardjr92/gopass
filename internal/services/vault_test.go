package services

import (
	"context"
	"errors"
	"testing"

	"github.com/edgardjr92/gopass/internal/keys"
	"github.com/edgardjr92/gopass/internal/mocks"
	"github.com/edgardjr92/gopass/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestNewVaultService(t *testing.T) {
	repoMock := &mocks.VaultRepositoryMock{}

	vaultSvc := NewVaultService(repoMock)

	assert.Equal(t, repoMock, vaultSvc.repository)
}

func TestCreateVault(t *testing.T) {
	userID := uint(10)
	ctx := context.WithValue(context.TODO(), keys.UserIDKey, userID)
	name := "My Vault"

	t.Run("success", func(t *testing.T) {
		// given
		repoMock := &mocks.VaultRepositoryMock{}

		repoMock.On("FindByNameAndUserID", ctx, name, userID).Return(&models.Vault{}, nil)
		repoMock.On("Save", ctx, &models.Vault{Name: name, UserID: userID}).Run(func(args mock.Arguments) {
			vault := args.Get(1).(*models.Vault)
			vault.ID = uint(100)
		})

		// when
		vaultSvc := &vaultService{repository: repoMock}
		actual, error := vaultSvc.Create(ctx, name)

		// then
		assert.Equal(t, uint(100), actual)
		assert.Nil(t, error)

		repoMock.AssertExpectations(t)
	})

	t.Run("user not authenticated", func(t *testing.T) {
		// given
		repoMock := &mocks.VaultRepositoryMock{}
		ctx := context.TODO()

		// when
		vaultSvc := &vaultService{repository: repoMock}
		actual, error := vaultSvc.Create(ctx, name)

		// then
		assert.Equal(t, uint(0), actual)
		assert.Equal(t, "UNAUTHORIZED: User is not authenticated", error.Error())

		repoMock.AssertExpectations(t)
	})

	names := []string{"", " "}
	for _, n := range names {
		t.Run("name is required", func(t *testing.T) {
			// given
			repoMock := &mocks.VaultRepositoryMock{}

			// when
			vaultSvc := &vaultService{repository: repoMock}
			actual, error := vaultSvc.Create(ctx, n)

			// then
			assert.Equal(t, uint(0), actual)
			assert.Equal(t, "BAD_REQUEST: name is required", error.Error())

			repoMock.AssertExpectations(t)
		})
	}

	t.Run("vault already exists", func(t *testing.T) {
		// given
		repoMock := &mocks.VaultRepositoryMock{}

		repoMock.On("FindByNameAndUserID", ctx, name, userID).
			Return(&models.Vault{Model: gorm.Model{ID: 1}}, nil)

		// when
		vaultSvc := &vaultService{repository: repoMock}
		actual, error := vaultSvc.Create(ctx, name)

		// then
		assert.Equal(t, uint(0), actual)
		assert.Equal(t, "RESOURCE_ALREADY_EXISTS: Vault already exists", error.Error())

		repoMock.AssertExpectations(t)
	})

	testCases := []struct {
		findByNameAndUserIDError error
		saveError                error
		expectedError            error
	}{
		{
			findByNameAndUserIDError: errors.New("error when finding vault"),
			saveError:                nil,
			expectedError:            errors.New("error when finding vault"),
		},
		{
			findByNameAndUserIDError: nil,
			saveError:                errors.New("error when saving vault"),
			expectedError:            errors.New("error when saving vault"),
		},
	}
	for _, tc := range testCases {
		t.Run("unexpected error", func(t *testing.T) {
			// given
			repoMock := &mocks.VaultRepositoryMock{}

			repoMock.On("FindByNameAndUserID", ctx, name, userID).Return(&models.Vault{}, tc.findByNameAndUserIDError)
			repoMock.On("Save", ctx, mock.Anything).Return(tc.saveError)

			// when
			vaultSvc := &vaultService{repository: repoMock}
			actual, error := vaultSvc.Create(ctx, name)

			// then
			assert.Equal(t, uint(0), actual)
			assert.Equal(t, tc.expectedError, error)
		})
	}

}

func TestGetAllVaults(t *testing.T) {
	userID := uint(10)
	ctx := context.WithValue(context.TODO(), keys.UserIDKey, userID)

	testCase := []struct {
		mockReturn []models.Vault
		expected   []models.VaultDetail
	}{
		{
			mockReturn: []models.Vault{},
			expected:   []models.VaultDetail{},
		},
		{
			mockReturn: []models.Vault{
				{Model: gorm.Model{ID: 1}, Name: "My Vault"},
				{Model: gorm.Model{ID: 2}, Name: "My Vault 2"},
			},
			expected: []models.VaultDetail{
				{ID: uint(1), Name: "My Vault"},
				{ID: uint(2), Name: "My Vault 2"},
			},
		},
	}
	for _, tc := range testCase {
		t.Run("success", func(t *testing.T) {
			// given
			repoMock := &mocks.VaultRepositoryMock{}

			repoMock.On("FindByUserID", ctx, userID).Return(tc.mockReturn, nil)

			// when
			vaultSvc := &vaultService{repository: repoMock}
			actual, error := vaultSvc.GetAll(ctx)

			// then
			assert.Equal(t, tc.expected, actual)
			assert.Nil(t, error)

			repoMock.AssertExpectations(t)
		})
	}

	t.Run("user not authenticated", func(t *testing.T) {
		// given
		repoMock := &mocks.VaultRepositoryMock{}
		ctx := context.TODO()

		// when
		vaultSvc := &vaultService{repository: repoMock}
		actual, error := vaultSvc.GetAll(ctx)

		// then
		assert.Equal(t, []models.VaultDetail{}, actual)
		assert.Equal(t, "UNAUTHORIZED: User is not authenticated", error.Error())

		repoMock.AssertExpectations(t)
	})

	t.Run("unexpected error", func(t *testing.T) {
		// given
		repoMock := &mocks.VaultRepositoryMock{}

		repoMock.On("FindByUserID", ctx, userID).Return([]models.Vault{}, errors.New("error when finding vaults"))

		// when
		vaultSvc := &vaultService{repository: repoMock}
		actual, error := vaultSvc.GetAll(ctx)

		// then
		assert.Equal(t, []models.VaultDetail{}, actual)
		assert.Equal(t, "error when finding vaults", error.Error())

		repoMock.AssertExpectations(t)
	})
}
