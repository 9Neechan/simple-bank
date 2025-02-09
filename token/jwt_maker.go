package token

import (
	"fmt"
	"time"
)

type JWTMaker struct {
	secretKey string
}

const minSecretKeySize = 32

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {return "", nil}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {return &Payload{}, nil}

// go get github.com/dgrijalva/jwt-go

// * можно реализовать JWT токен при желании... но его пока что нет.... (20)
