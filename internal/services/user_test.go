package services

import (
	"context"
	"errors"
	"testing"

	cerrors "github.com/edgardjr92/gopass/internal/errors"
	"github.com/edgardjr92/gopass/internal/mocks"
	"github.com/edgardjr92/gopass/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	ctx := context.TODO()
	name := "John Doe"
	email := "jhon@test.com"
	psw := "123456"
	confirmPsw := "123456"

	t.Run("success", func(t *testing.T) {
		// given
		repoMock := &mocks.UserRepositoryMock{}
		hasherMock := &mocks.HasherMock{}

		repoMock.On("FindByEmail", email).Return(&models.User{}, nil)
		hasherMock.On("HashPsw", psw).Return("hashed-password", nil)
		repoMock.On(
			"Save",
			&models.User{Name: name, Email: email, Psw: "hashed-password"},
		).Return(1, nil)

		// when
		userSvc := &userService{
			repository: repoMock,
			hasher:     hasherMock,
		}
		actual, error := userSvc.Create(ctx, name, email, psw, confirmPsw)

		// then
		assert.Equal(t, 1, actual)
		assert.Nil(t, error)

		repoMock.AssertExpectations(t)
		hasherMock.AssertExpectations(t)
	})

	t.Run("passwords do not match", func(t *testing.T) {
		// given
		repoMock := &mocks.UserRepositoryMock{}
		hasherMock := &mocks.HasherMock{}

		// when
		userSvc := &userService{
			repository: repoMock,
			hasher:     hasherMock,
		}
		actual, error := userSvc.Create(ctx, name, email, psw, "1234567")

		// then
		assert.Equal(t, 0, actual)
		assert.Equal(t, error, cerrors.UnprocessableError("Passwords do not match"))
	})

	t.Run("user already exists", func(t *testing.T) {
		// given
		repoMock := &mocks.UserRepositoryMock{}
		hasherMock := &mocks.HasherMock{}

		repoMock.On("FindByEmail", email).
			Return(&models.User{Model: gorm.Model{ID: 1}}, nil)

		// when
		userSvc := &userService{
			repository: repoMock,
			hasher:     hasherMock,
		}
		actual, error := userSvc.Create(ctx, name, email, psw, confirmPsw)

		// then
		assert.Equal(t, 0, actual)
		assert.Equal(t, error, cerrors.ConflictError("User already exists"))
	})

	testCases := []struct {
		findError     error
		hashError     error
		saveError     error
		expectedError error
	}{
		{
			findError:     errors.New("error when finding user"),
			hashError:     nil,
			saveError:     nil,
			expectedError: errors.New("error when finding user"),
		},
		{
			findError:     nil,
			hashError:     errors.New("error when hashing password"),
			saveError:     nil,
			expectedError: errors.New("error when hashing password"),
		},
		{
			findError:     nil,
			hashError:     nil,
			saveError:     errors.New("error when saving user"),
			expectedError: errors.New("error when saving user"),
		},
	}

	for _, tc := range testCases {
		t.Run("Unexpected error", func(t *testing.T) {
			// given
			repoMock := &mocks.UserRepositoryMock{}
			hasherMock := &mocks.HasherMock{}

			repoMock.On("FindByEmail", email).Return(&models.User{}, tc.findError)
			hasherMock.On("HashPsw", psw).Return("", tc.hashError)
			repoMock.On("Save", mock.Anything).Return(0, tc.saveError)

			// when
			userSvc := &userService{
				repository: repoMock,
				hasher:     hasherMock,
			}
			actual, error := userSvc.Create(ctx, name, email, psw, confirmPsw)

			// then
			assert.Equal(t, 0, actual)
			assert.Equal(t, error, tc.expectedError)
		})
	}

}
