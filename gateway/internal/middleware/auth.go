package middleware

import (
	"context"
	"net/http"
	"os"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
)

type ContextString string

const ClerkClaimsKey ContextString = "clerk-claims"

const DevUserKey ContextString = "dev-user"

func AuthMiddleware(next http.Handler, clerkKey string) http.Handler {
	clerk.SetKey(os.Getenv(clerkKey))

	return clerkhttp.WithHeaderAuthorization()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("ENV") == "dev" {
			// skip auth in dev mode
			ctx := context.WithValue(r.Context(), DevUserKey, "dev-user")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		claims, ok := clerk.SessionClaimsFromContext(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ClerkClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}))
}
