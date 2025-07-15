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

func GenerateTokenPair(userID string) (accessToken string, refreshToken string, refreshExp time.Time, err error) {
	accessToken, err = GenerateToken(userID, 15*time.Minute)
	if err != nil {
		return "", "", time.Time{}, err
	}

	refreshTTL := 7 * 24 * time.Hour
	refreshExp = time.Now().Add(refreshTTL)

	refreshToken, err = GenerateToken(userID, refreshTTL)
	if err != nil {
		return "", "", time.Time{}, err
	}

	return accessToken, refreshToken, refreshExp, nil
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
}
