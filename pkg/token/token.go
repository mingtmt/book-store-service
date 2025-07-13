package token

import (
	"crypto/rsa"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

func InitKeys() error {
	keyPath := os.Getenv("KEY_PATH")
	privData, err := os.ReadFile(keyPath + "/private.pem")
	if err != nil {
		return err
	}
	pubData, err := os.ReadFile(keyPath + "/public.pem")
	if err != nil {
		return err
	}

	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privData)
	if err != nil {
		return err
	}
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(pubData)
	return err
}

func GenerateToken(userID string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(duration).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

func GenerateTokenPair(userID string) (accessToken, refreshToken string, err error) {
	accessToken, err = GenerateToken(userID, 15*time.Minute) // short-lived
	if err != nil {
		return "", "", err
	}
	refreshToken, err = GenerateToken(userID, 7*24*time.Hour) // long-lived
	return
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
}
