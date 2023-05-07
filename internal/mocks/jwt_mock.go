package mocks

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type JWTGeneratorMock struct {
	mock.Mock
}

func (m *JWTGeneratorMock) Generate(userID uint, exp time.Time) (string, error) {
	args := m.Called(userID, exp)
	return args.String(0), args.Error(1)
}
