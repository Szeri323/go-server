package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Password auth
func HashPassword(password string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("couldn't hash the password: %w", err)
	}

	return string(hashedPassword), nil

}

func CheckPasswordHash(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// Token auth
func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	})
	signedToken, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("couldn't signed the token: %w", err)
	}

	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	clamis := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &clamis, func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil })
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("couldn't parse the token: %w", err)
	}
	if !token.Valid {
		return uuid.UUID{}, fmt.Errorf("token is not valid: %w", err)
	}
	id, err := uuid.Parse(clamis.Subject)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("couldn't parse the uuid: %w", err)
	}
	return id, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	for _, header := range headers {
		str := "Bearer "
		if strings.Contains(header[0], str) {
			newstring, check := strings.CutPrefix(header[0], str)
			if check {
				return newstring, nil
			}
		}
	}
	return "", fmt.Errorf("header doesn't exists")
}
func GetAPIKey(headers http.Header) (string, error) {
	for _, header := range headers {
		str := "ApiKey "
		if strings.Contains(header[0], str) {
			newstring, check := strings.CutPrefix(header[0], str)
			if check {
				return newstring, nil
			}
		}
	}
	return "", fmt.Errorf("header doesn't exists")
}

func MakeRefreshToken() (string, error) {
	randByte := make([]byte, 32)
	_, err := rand.Read(randByte)
	if err != nil {
		return "", fmt.Errorf("couldn't populete byte arr: %w", err)
	}
	str := hex.EncodeToString(randByte)
	return str, nil
}
