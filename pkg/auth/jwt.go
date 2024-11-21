package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	privateKey []byte
	publicKey  []byte
)

func init() {
	var err error

	// Load the private key from the file
	privateKey, err = os.ReadFile("certs/jwt/private.pem")
	if err != nil {
		panic("unable to read private key: " + err.Error())
	}

	// Load the public key from the file
	publicKey, err = os.ReadFile("certs/jwt/public.pem")
	if err != nil {
		panic("unable to read public key: " + err.Error())
	}
}

type Claims struct {
	Username  string `json:"username"`
	Phone     string `json:"phone"`
	Notes     string `json:"notes"`
	ExpiresIn int    `json:"expires_in"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a JWT token for a given username
func GenerateJWT(username, phone, notes string, expiresIn int) (string, error) {
	if expiresIn == 0 {
		expiresIn = 30 // Default expiration of 30 days
	}

	expirationTime := time.Now().Add(time.Duration(expiresIn) * 24 * time.Hour)
	claims := &Claims{
		Username:  username,
		Phone:     phone,
		Notes:     notes,
		ExpiresIn: expiresIn,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Parse the private key for signing the token
	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ValidateJWT validates the JWT token and returns the claims
func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	// Parse the public key for validating the token
	key, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, errors.New("invalid signature")
		}
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	if claims.RegisteredClaims.ExpiresAt != nil && time.Now().After(claims.RegisteredClaims.ExpiresAt.Time) {
		return nil, errors.New("token has expired")
	}
	return claims, nil
}
