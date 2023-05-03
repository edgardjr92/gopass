package hash

import "golang.org/x/crypto/bcrypt"

type Hasher interface {
	HashPsw(password string) (string, error)
}

type bcryptHasher struct{}

func NewBcryptHasher() *bcryptHasher {
	return &bcryptHasher{}
}

func (b *bcryptHasher) HashPsw(password string) (string, error) {
	// Generate a salt for the hash
	salt, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	// Generate the hash using the password and salt
	hash, err := bcrypt.GenerateFromPassword(salt, bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}
