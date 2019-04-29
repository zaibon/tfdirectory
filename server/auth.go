package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

const iyoBaseURL = "https://itsyou.online/v1"

var iyoKey = []byte(`-----BEGIN PUBLIC KEY-----
MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAES5X8XrfKdx9gYayFITc89wad4usrk0n2
7MjiGYvqalizeSWTHEpnd7oea9IQ8T5oJjMVH5cc0H5tFSKilFFeh//wngxIyny6
6+Vq5t5B0V0Ehy01+2ceEon2Y0XDkIKv
-----END PUBLIC KEY-----`)

type contextOrgType struct{}

var contextOrgKey = &contextOrgType{}

type ScopeClaim struct {
	Scope []string `json:"scope"`
	jwt.StandardClaims
}

// ExtractScopeMiddleware extract the scope from any jwt token passed in
// the Authorization header and put it in the request context
func ExtractScopeMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			next.ServeHTTP(w, r)
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &ScopeClaim{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return jwt.ParseECPublicKeyFromPEM(iyoKey)
		})
		if err != nil {
			log.Printf("fail to parse jwt token: %v", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*ScopeClaim)
		if !ok || !token.Valid {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if len(claims.Scope) < 1 {
			http.Error(w, "missing scope claim", http.StatusUnauthorized)
			return
		}

		log.Printf("%+v", claims.Scope)
		// scopesStr, ok := scopes.([]string)
		// if !ok {
		// 	http.Error(w, "wrong type for scope claim", http.StatusUnauthorized)
		// 	return
		// }

		// progagate the scopes in the request context
		r = r.WithContext(withScope(r.Context(), claims.Scope))

		next.ServeHTTP(w, r)
	})
}

func withScope(ctx context.Context, scope []string) context.Context {
	scopeStr := strings.Join(scope, ",")
	return context.WithValue(ctx, contextOrgKey, scopeStr)
}

// ScopeFromContext returns the scope in the jwt from the context.
// A zero ID is returned if there are no identifers in the
// current context.
func ScopeFromContext(ctx context.Context) []string {
	v := ctx.Value(contextOrgKey)
	if v == nil {
		return []string{""}
	}
	vStr := v.(string)
	return strings.Split(vStr, ",")
}
