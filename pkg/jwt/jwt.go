package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTGenerator interface {
	Generate(userID uint, exp time.Time) (string, error)
}

type jwtGo struct {
	secret []byte
}

func NewJWTService(secret []byte) JWTGenerator {
	return &jwtGo{secret: secret}
}

// Generate generates a new JWT token for a given user ID and expiration time.
// The JWT token is signed using HMAC-SHA256 and a secret key. The resulting
// token string can be used to authenticate the user in subsequent requests.
//
// userID: the user ID to include in the JWT token claims.
// exp: the expiration time of the JWT token.
//
// Returns the JWT token string or an error if the token could not be generated.
func (j jwtGo) Generate(userID uint, exp time.Time) (string, error) {
	if exp.IsZero() {
		exp = time.Now().Add(24 * time.Hour)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     exp.Unix(),
	})

	tokenStr, err := token.SignedString(j.secret)

	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
