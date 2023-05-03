package services

import (
	"context"
	"errors"
	"testing"

	cerrors "github.com/edgardjr92/gopass/internal/errors"
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
		mockRepo := &mockRepository{}
		mockHasher := &mockHasher{}

		mockRepo.On("FindByEmail", email).Return(&models.User{}, nil)
		mockHasher.On("HashPsw", psw).Return("hashed-password", nil)
		mockRepo.On(
			"Save",
			&models.User{Name: name, Email: email, Psw: "hashed-password"},
		).Return(1, nil)

		// when
		userSvc := &userService{
			repository: mockRepo,
			hasher:     mockHasher,
		}
		actual, error := userSvc.Create(ctx, name, email, psw, confirmPsw)

		// then
		assert.Equal(t, 1, actual)
		assert.Nil(t, error)

		mockRepo.AssertExpectations(t)
		mockHasher.AssertExpectations(t)
	})

	t.Run("passwords do not match", func(t *testing.T) {
		// given
		mockRepo := &mockRepository{}
		mockHasher := &mockHasher{}

		// when
		userSvc := &userService{
			repository: mockRepo,
			hasher:     mockHasher,
		}
		actual, error := userSvc.Create(ctx, name, email, psw, "1234567")

		// then
		assert.Equal(t, 0, actual)
		assert.Equal(t, error, cerrors.UnprocessableError("Passwords do not match"))
	})

	t.Run("user already exists", func(t *testing.T) {
		// given
		mockRepo := &mockRepository{}
		mockHasher := &mockHasher{}

		mockRepo.On("FindByEmail", email).
			Return(&models.User{Model: gorm.Model{ID: 1}}, nil)

		// when
		userSvc := &userService{
			repository: mockRepo,
			hasher:     mockHasher,
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
			mockRepo := &mockRepository{}
			mockHasher := &mockHasher{}

			mockRepo.On("FindByEmail", email).Return(&models.User{}, tc.findError)
			mockHasher.On("HashPsw", psw).Return("", tc.hashError)
			mockRepo.On("Save", mock.Anything).Return(0, tc.saveError)

			// when
			userSvc := &userService{
				repository: mockRepo,
				hasher:     mockHasher,
			}
			actual, error := userSvc.Create(ctx, name, email, psw, confirmPsw)

			// then
			assert.Equal(t, 0, actual)
			assert.Equal(t, error, tc.expectedError)
		})
	}

}

// Define mock repository
type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockRepository) Save(user *models.User) (int, error) {
	args := m.Called(user)
	return args.Int(0), args.Error(1)
}

// Define mock hasher
type mockHasher struct {
	mock.Mock
}

func (m *mockHasher) HashPsw(psw string) (string, error) {
	args := m.Called(psw)
	return args.String(0), args.Error(1)
}
