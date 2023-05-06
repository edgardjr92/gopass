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

func TestCreateVault(t *testing.T) {
	userID := uint(10)
	cxt := context.WithValue(context.TODO(), keys.UserIDKey, userID)
	name := "My Vault"

	t.Run("success", func(t *testing.T) {
		// given
		repoMock := &mocks.VaultRepositoryMock{}

		repoMock.On("FindByNameAndUserID", name, userID).Return(&models.Vault{}, nil)
		repoMock.On("Save", &models.Vault{Name: name, UserID: userID}).Run(func(args mock.Arguments) {
			vault := args.Get(0).(*models.Vault)
			vault.ID = uint(100)
		})

		// when
		vaultSvc := &vaultService{repository: repoMock}
		actual, error := vaultSvc.Create(cxt, name)

		// then
		assert.Equal(t, uint(100), actual)
		assert.Nil(t, error)

		repoMock.AssertExpectations(t)
	})

	t.Run("user not authenticated", func(t *testing.T) {
		// given
		repoMock := &mocks.VaultRepositoryMock{}
		cxt := context.TODO()

		// when
		vaultSvc := &vaultService{repository: repoMock}
		actual, error := vaultSvc.Create(cxt, name)

		// then
		assert.Equal(t, uint(0), actual)
		assert.Equal(t, "UNAUTHORIZED_ERROR: User is not authenticated", error.Error())

		repoMock.AssertExpectations(t)
	})

	names := []string{"", " "}
	for _, n := range names {
		t.Run("name is required", func(t *testing.T) {
			// given
			repoMock := &mocks.VaultRepositoryMock{}

			// when
			vaultSvc := &vaultService{repository: repoMock}
			actual, error := vaultSvc.Create(cxt, n)

			// then
			assert.Equal(t, uint(0), actual)
			assert.Equal(t, "BAD_REQUEST_ERROR: name is required", error.Error())

			repoMock.AssertExpectations(t)
		})
	}

	t.Run("vault already exists", func(t *testing.T) {
		// given
		repoMock := &mocks.VaultRepositoryMock{}

		repoMock.On("FindByNameAndUserID", name, userID).
			Return(&models.Vault{Model: gorm.Model{ID: 1}}, nil)

		// when
		vaultSvc := &vaultService{repository: repoMock}
		actual, error := vaultSvc.Create(cxt, name)

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

			repoMock.On("FindByNameAndUserID", name, userID).Return(&models.Vault{}, tc.findByNameAndUserIDError)
			repoMock.On("Save", mock.Anything).Return(tc.saveError)

			// when
			vaultSvc := &vaultService{repository: repoMock}
			actual, error := vaultSvc.Create(cxt, name)

			// then
			assert.Equal(t, uint(0), actual)
			assert.Equal(t, tc.expectedError, error)
		})
	}

}
