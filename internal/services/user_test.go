package services

import (
	"context"
	"errors"
	"testing"

	"github.com/edgardjr92/gopass/internal/cerrors"
	"github.com/edgardjr92/gopass/internal/mocks"
	"github.com/edgardjr92/gopass/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestNewUserService(t *testing.T) {
	repoMock := &mocks.UserRepositoryMock{}
	hasherMock := &mocks.HasherMock{}

	userSvc := NewUserService(repoMock, hasherMock)

	assert.Equal(t, repoMock, userSvc.repository)
	assert.Equal(t, hasherMock, userSvc.hasher)
}

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

		newUser := &models.User{Name: name, Email: email, Psw: "hashed-password"}

		repoMock.On("FindByEmail", ctx, email).Return(&models.User{}, nil)
		hasherMock.On("HashPsw", psw).Return("hashed-password", nil)
		repoMock.On("Save", ctx, newUser).Run(func(args mock.Arguments) {
			user := args.Get(1).(*models.User)
			user.ID = uint(1)
		})

		// when
		userSvc := &userService{repository: repoMock, hasher: hasherMock}
		actual, error := userSvc.Create(ctx, name, email, psw, confirmPsw)

		// then
		assert.Equal(t, uint(1), actual)
		assert.Nil(t, error)

		repoMock.AssertExpectations(t)
		hasherMock.AssertExpectations(t)
	})

	t.Run("passwords do not match", func(t *testing.T) {
		// given
		repoMock := &mocks.UserRepositoryMock{}
		hasherMock := &mocks.HasherMock{}

		// when
		userSvc := &userService{repository: repoMock, hasher: hasherMock}
		actual, error := userSvc.Create(ctx, name, email, psw, "1234567")

		// then
		assert.Equal(t, uint(0), actual)
		assert.Equal(t, error, cerrors.UnprocessableError("passwords do not match"))
	})

	t.Run("user already exists", func(t *testing.T) {
		// given
		repoMock := &mocks.UserRepositoryMock{}
		hasherMock := &mocks.HasherMock{}

		repoMock.On("FindByEmail", ctx, email).
			Return(&models.User{Model: gorm.Model{ID: 1}}, nil)

		// when
		userSvc := &userService{repository: repoMock, hasher: hasherMock}
		actual, error := userSvc.Create(ctx, name, email, psw, confirmPsw)

		// then
		assert.Equal(t, uint(0), actual)
		assert.Equal(t, error, cerrors.ConflictError("user already exists"))
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
		t.Run("unexpected error", func(t *testing.T) {
			// given
			repoMock := &mocks.UserRepositoryMock{}
			hasherMock := &mocks.HasherMock{}

			repoMock.On("FindByEmail", ctx, email).Return(&models.User{}, tc.findError)
			hasherMock.On("HashPsw", psw).Return("", tc.hashError)
			repoMock.On("Save", ctx, mock.Anything).Return(tc.saveError)

			// when
			userSvc := &userService{repository: repoMock, hasher: hasherMock}
			actual, error := userSvc.Create(ctx, name, email, psw, confirmPsw)

			// then
			assert.Equal(t, uint(0), actual)
			assert.Equal(t, tc.expectedError, error)
		})
	}

}
