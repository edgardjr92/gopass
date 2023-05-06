package mocks

import (
	"github.com/edgardjr92/gopass/internal/models"
	"github.com/stretchr/testify/mock"
)

// Define mock repository
type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserRepositoryMock) Save(user *models.User) error {
	args := m.Called(user)
	if len(args) > 0 {
		err, _ := args[0].(error)
		return err
	}
	return nil
}
