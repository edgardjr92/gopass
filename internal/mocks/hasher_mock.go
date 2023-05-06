package mocks

import "github.com/stretchr/testify/mock"

// Define mock hasher
type HasherMock struct {
	mock.Mock
}

func (m *HasherMock) HashPsw(psw string) (string, error) {
	args := m.Called(psw)
	return args.String(0), args.Error(1)
}
