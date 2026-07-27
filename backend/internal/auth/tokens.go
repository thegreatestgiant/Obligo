package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func Verifyer(tokenString string, secret []byte) (claims *jwt.RegisteredClaims, err error) {
	c := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, c, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		log.Printf("Cannot Parse Token: %v", err)
		return &jwt.RegisteredClaims{}, err
	}
	if !token.Valid {
		log.Println("Invalid token:")
		return &jwt.RegisteredClaims{}, errors.New("Invalid Token")
	}
	return c, nil
}

func VerifyerOfExpired(tokenString string, secret []byte) (claims *jwt.RegisteredClaims, err error) {
	c := &jwt.RegisteredClaims{}
	_, err = jwt.ParseWithClaims(tokenString, c, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		log.Printf("Cannot Parse Token (fatal): %v", err)
		return &jwt.RegisteredClaims{}, err
	}

	return c, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret []byte, expiresIn time.Duration) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "Yoily",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
		ID:        uuid.NewString(),
	}).SignedString(tokenSecret)
}

func MakeRefreshToken() string {
	token := make([]byte, 32)
	rand.Read(token)
	return hex.EncodeToString(token)
}
