package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/edgardjr92/gopass/internal/mocks"
	"github.com/edgardjr92/gopass/internal/models"
	"github.com/edgardjr92/gopass/pkg/clock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestNewAuthService(t *testing.T) {
	jwtMock := &mocks.JWTGeneratorMock{}
	repoMock := &mocks.UserRepositoryMock{}
	clockMock := clock.Clock{}

	authSrv := NewAuthService(jwtMock, repoMock, clockMock)

	assert.NotNil(t, authSrv)
	assert.Equal(t, jwtMock, authSrv.jwt)
	assert.Equal(t, repoMock, authSrv.repository)
	assert.Equal(t, clockMock, authSrv.clock)
}

func TestLogin(t *testing.T) {
	ctx := context.TODO()
	email := "test@test.com"
	authKey := "hashed-auth-key"

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"

	clockMock := clock.Clock{
		NowFn: func() time.Time {
			return time.Date(2023, 5, 6, 0, 0, 0, 0, time.UTC)
		},
	}

	t.Run("success", func(t *testing.T) {
		// given
		jwtMock := &mocks.JWTGeneratorMock{}
		repoMock := &mocks.UserRepositoryMock{}

		repoMock.On("FindByEmail", ctx, email).
			Return(&models.User{Model: gorm.Model{ID: 1}, AuthKey: authKey}, nil)

		tomorrow := time.Date(2023, 5, 7, 0, 0, 0, 0, time.UTC)
		jwtMock.On("Generate", uint(1), tomorrow).Return(token, nil)

		// when
		authSrv := &authService{jwtMock, clockMock, repoMock}
		actual, error := authSrv.Login(ctx, email, authKey)

		// then
		assert.Equal(t, token, actual)
		assert.Nil(t, error)

		repoMock.AssertExpectations(t)
		jwtMock.AssertExpectations(t)
	})

	args := []struct {
		name    string
		email   string
		authKey string
		err     string
	}{
		{"empty email", "", "hashed-auth-key", "email is required"},
		{"blank email", "   ", "hashed-auth-key", "email is required"},
		{"empty password", "test@test.com", "", "authKey is required"},
		{"blank password", "test@test.com", "   ", "authKey is required"},
	}
	for _, arg := range args {
		t.Run(arg.name, func(t *testing.T) {
			// given
			jwtMock := &mocks.JWTGeneratorMock{}
			repoMock := &mocks.UserRepositoryMock{}

			// when
			authSrv := &authService{jwtMock, clockMock, repoMock}
			actual, error := authSrv.Login(ctx, arg.email, arg.authKey)

			// then
			assert.Equal(t, "", actual)
			assert.Equal(t, arg.err, error.Error())
		})
	}

	t.Run("user not found", func(t *testing.T) {
		// given
		jwtMock := &mocks.JWTGeneratorMock{}
		repoMock := &mocks.UserRepositoryMock{}

		repoMock.On("FindByEmail", ctx, email).Return(&models.User{}, nil)

		// when
		authSrv := &authService{jwtMock, clockMock, repoMock}
		actual, error := authSrv.Login(ctx, email, authKey)

		// then
		assert.Equal(t, "", actual)
		assert.Equal(t, "invalid credentials", error.Error())
	})

	t.Run("invalid authKey", func(t *testing.T) {
		// given
		jwtMock := &mocks.JWTGeneratorMock{}
		repoMock := &mocks.UserRepositoryMock{}

		repoMock.On("FindByEmail", ctx, email).
			Return(&models.User{Model: gorm.Model{ID: 1}, AuthKey: authKey}, nil)

		// when
		authSrv := &authService{jwtMock, clockMock, repoMock}
		actual, error := authSrv.Login(ctx, email, "invalid-auth-key")

		// then
		assert.Equal(t, "", actual)
		assert.Equal(t, "invalid credentials", error.Error())
	})

	testCases := []struct {
		name             string
		findByEmailError error
		generateError    error
		expected         error
	}{
		{
			name:             "userRepository.FindByEmail error",
			findByEmailError: fmt.Errorf("error finding user"),
			generateError:    nil,
			expected:         fmt.Errorf("error finding user"),
		},
		{
			name:             "jwtGenerator.Generate error",
			findByEmailError: nil,
			generateError:    fmt.Errorf("error generating token"),
			expected:         fmt.Errorf("error generating token"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			jwtMock := &mocks.JWTGeneratorMock{}
			repoMock := &mocks.UserRepositoryMock{}

			repoMock.On("FindByEmail", ctx, email).
				Return(&models.User{Model: gorm.Model{ID: 1}, AuthKey: authKey}, tc.findByEmailError)

			jwtMock.On("Generate", uint(1), mock.Anything).Return("", tc.generateError)

			// when
			authSrv := &authService{jwtMock, clockMock, repoMock}
			actual, error := authSrv.Login(ctx, email, authKey)

			// then
			assert.Equal(t, "", actual)
			assert.Equal(t, tc.expected, error)
		})
	}

}
