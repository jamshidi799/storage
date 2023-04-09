package user

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"storage/domain"
	"strconv"
)

type jwtTokenGenerator struct {
	secret []byte
}

func NewJwtTokenGenerator(secret string) domain.TokenGenerator {
	return &jwtTokenGenerator{secret: []byte(secret)}
}

func (j *jwtTokenGenerator) Verify(tokenStr string) bool {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})
	if err != nil {
		return false
	}
	_, ok := token.Claims.(jwt.MapClaims)
	return ok && token.Valid

}

func (j *jwtTokenGenerator) Generate(id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": strconv.Itoa(id),
	})

	return token.SignedString(j.secret)
}
