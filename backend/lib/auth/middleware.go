package auth

import (
	"context"
	"doodemo-cook/lib/response"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
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

func Middleware() func(next http.Handler) http.Handler {
	domain := os.Getenv("AUTH0_DOMAIN")
	if domain == "" {
		log.Fatal("you must set your 'AUTH0_DOMAIN' environmental variable")
	}
	audience := os.Getenv("AUTH0_AUDIENCE")
	if audience == "" {
		log.Fatal("you must set your 'AUTH0_AUDIENCE' environmental variable")
	}

	issuerURL, err := url.Parse("https://" + domain + "/")
	if err != nil {
		log.Fatalf("failed to parse the issuer url: %v", err)
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{audience},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		log.Fatalf("failed to set up the jwt validator")
	}
	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("encountered error while validating JWT: %v", err)
		response.FromStatusCode(r.Context(), w, http.StatusUnauthorized)
	}

	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
	)

	return func(next http.Handler) http.Handler {
		return middleware.CheckJWT(next)
	}
}
