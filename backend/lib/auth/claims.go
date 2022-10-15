package auth

import (
	"context"
	"errors"
	"strings"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
)

type CustomClaims struct {
	Scope string `json:"scope"`
}

func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}

func (c CustomClaims) HasScope(expectedScope string) bool {
	result := strings.Split(c.Scope, " ")
	for i := range result {
		if result[i] == expectedScope {
			return true
		}
	}
	return false
}

func ClaimsFromContext(ctx context.Context) (*CustomClaims, error) {
	token, ok := ctx.Value(jwtmiddleware.ContextKey{}).(validator.ValidatedClaims)
	if !ok {
		return nil, errors.New("context doesn't have jwt token")
	}

	claims, ok := token.CustomClaims.(*CustomClaims)
	if !ok {
		return nil, errors.New("claims is not CustomClaims")
	}

	return claims, nil
}
