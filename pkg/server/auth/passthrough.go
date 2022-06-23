package auth

import (
	"fmt"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-logr/logr"
)

// NewJWTPassthroughCookiePrincipalGetter creates and returns a new
// JWTPassthroughCookiePrincipalGetter.
func NewJWTPassthroughCookiePrincipalGetter(log logr.Logger, verifier *oidc.IDTokenVerifier, cookieName string) PrincipalGetter {
	return &JWTPassthroughCookiePrincipalGetter{
		log:        log,
		verifier:   verifier,
		cookieName: cookieName,
	}
}

// JWTPassthroughCookiePrincipalGetter inspects a cookie for a JWT token and returns a
// principal value.
type JWTPassthroughCookiePrincipalGetter struct {
	log        logr.Logger
	verifier   *oidc.IDTokenVerifier
	cookieName string
}

func (pg *JWTPassthroughCookiePrincipalGetter) Principal(r *http.Request) (*UserPrincipal, error) {
	cookie, err := r.Cookie(pg.cookieName)
	if err == http.ErrNoCookie {
		return nil, nil
	}

	principal, err := parseJWTToken(r.Context(), pg.verifier, cookie.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to parse for passthrough: %w", err)
	}
	principal.Token = cookie.Value

	return principal, nil
}
