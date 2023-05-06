package services

import (
	"context"
	"log"

	"github.com/edgardjr92/gopass/internal/errors"
	"github.com/edgardjr92/gopass/internal/models"
	"github.com/edgardjr92/gopass/internal/repositories"
	"github.com/edgardjr92/gopass/pkg/hash"
)

type IUserService interface {
	// Create creates a new user.
	// It returns the ID of the newly created user.
	Create(ctx context.Context, name, email, psw, confirmPsw string) (uint, error)
}

type userService struct {
	repository repositories.IUserRepository
	hasher     hash.Hasher
}

func NewUserService(repository repositories.IUserRepository, hasher hash.Hasher) *userService {
	return &userService{repository, hasher}
}

func (u *userService) Create(ctx context.Context, name, email, psw, confirmPsw string) (uint, error) {
	if psw != confirmPsw {
		return 0, errors.UnprocessableError("Passwords do not match")
	}

	user, err := u.repository.FindByEmail(email)

	if err != nil {
		log.Printf("Error while trying to find user by email: %v", err.Error())
		return 0, err
	}

	if user.ID != 0 {
		return 0, errors.ConflictError("User already exists")
	}

	hashedPsw, err := u.hasher.HashPsw(psw)

	if err != nil {
		log.Printf("Error while trying to hash password: %v", err.Error())
		return 0, err
	}

	newUser := models.User{
		Name:  name,
		Email: email,
		Psw:   hashedPsw,
	}

	if err := u.repository.Save(&newUser); err != nil {
		log.Printf("Error while trying to save user: %v", err.Error())
		return 0, err
	}

	return newUser.ID, nil
}
