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
	clerk.SetKey(clerkKey)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("ENV") == "dev" {
			ctx := context.WithValue(r.Context(), DevUserKey, "dev-user")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		clerkhttp.WithHeaderAuthorization()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// If Clerk validates, attach claims
			if claims, ok := clerk.SessionClaimsFromContext(r.Context()); ok {
				ctx := context.WithValue(r.Context(), ClerkClaimsKey, claims)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// Otherwise just pass context through —
			// gqlgen @isPublic directive will decide
			next.ServeHTTP(w, r)
		})).ServeHTTP(w, r)
	})
}
