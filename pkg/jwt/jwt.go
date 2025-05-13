package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Config struct {
	SecretKey          []byte
	TokenExpiration    time.Duration
	Issuer             string
	SigningMethod      jwt.SigningMethod
	ValidateExpiration bool
	ValidateIssuer     bool
}

func DefaultConfig() *Config {
	return &Config{
		SecretKey:          []byte("default-secret-key"),
		TokenExpiration:    24 * time.Hour,
		Issuer:             "my-app",
		SigningMethod:      jwt.SigningMethodHS256,
		ValidateExpiration: true,
		ValidateIssuer:     false,
	}
}

type User struct {
	ID uuid.UUID `json:"id"`
}

type JWTManager struct {
	config *Config
}

type opt func(c *Config)

func WithSecretKey(key []byte) opt {
	return func(c *Config) {
		c.SecretKey = key
	}
}
func WithTokenExpiration(duration time.Duration) opt {
	return func(c *Config) {
		c.TokenExpiration = duration
	}
}

func NewJWTManager(params ...opt) *JWTManager {

	config := DefaultConfig()

	for _, param := range params {
		param(config)
	}

	return &JWTManager{config: config}
}

type UserClaims struct {
	User User `json:"user"`
	jwt.RegisteredClaims
}

func (m *JWTManager) GenerateToken(user User) (string, error) {
	claims := UserClaims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.config.TokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    m.config.Issuer,
		},
	}

	token := jwt.NewWithClaims(m.config.SigningMethod, claims)
	signedToken, err := token.SignedString(m.config.SecretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

func (m *JWTManager) ParseToken(tokenString string) (*User, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if token.Method != m.config.SigningMethod {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.config.SecretKey, nil
	}

	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, keyFunc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		if m.config.ValidateExpiration {
			if time.Now().After(claims.ExpiresAt.Time) {
				return nil, errors.New("token has expired")
			}
		}

		if m.config.ValidateIssuer && claims.Issuer != m.config.Issuer {
			return nil, errors.New("invalid token issuer")
		}

		return &claims.User, nil
	}

	return nil, errors.New("invalid token")
}
