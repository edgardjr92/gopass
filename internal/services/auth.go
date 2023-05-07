package services

import (
	"context"
	"log"
	"time"

	"github.com/edgardjr92/gopass/internal/cerrors"
	"github.com/edgardjr92/gopass/internal/repositories"
	"github.com/edgardjr92/gopass/internal/utils"
	"github.com/edgardjr92/gopass/pkg/clock"
	"github.com/edgardjr92/gopass/pkg/jwt"
)

type IAuthService interface {
	// Login authenticates a user.
	// It returns the JWT token string or an error if the user could not be authenticated.
	Login(ctx context.Context, email, authKey string) (string, error)
}

type authService struct {
	jwt        jwt.JWTGenerator
	clock      clock.Clock
	repository repositories.IUserRepository
}

func NewAuthService(jwt jwt.JWTGenerator, repository repositories.IUserRepository, clock clock.Clock) *authService {
	return &authService{jwt, clock, repository}
}

func (a *authService) Login(ctx context.Context, email, authKey string) (string, error) {
	if utils.IsBlank(email) {
		return "", cerrors.BadRequestError("email is required")
	}

	if utils.IsBlank(authKey) {
		return "", cerrors.BadRequestError("authKey is required")
	}

	user, err := a.repository.FindByEmail(ctx, email)

	if err != nil {
		log.Printf("error while trying to find user by email: %v", err.Error())
		return "", err
	}

	if user.ID == 0 || user.AuthKey != authKey {
		return "", cerrors.UnauthorizedError("invalid credentials")
	}

	token, err := a.jwt.Generate(user.ID, a.clock.Now().Add(24*time.Hour))

	if err != nil {
		log.Printf("error while trying to generate JWT token: %v", err.Error())
		return "", err
	}

	return token, nil
}
