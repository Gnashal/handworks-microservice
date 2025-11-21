package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/gin-gonic/gin"
)

type ContextString string

const ClerkClaimsKey ContextString = "clerk-claims"

func ClerkAuthMiddleware(publicPaths []string) gin.HandlerFunc {
	clerkKey := os.Getenv("CLERK_SECRET_KEY")
	clerk.SetKey(clerkKey)

	return func(c *gin.Context) {
		// Skip public paths
		for _, path := range publicPaths {
			if strings.HasPrefix(c.FullPath(), path) {
				c.Next()
				return
			}
		}
		clerkhttp.WithHeaderAuthorization()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if claims, ok := clerk.SessionClaimsFromContext(r.Context()); ok {
				c.Set(string(ClerkClaimsKey), claims)
				c.Next()
				return
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		})).ServeHTTP(c.Writer, c.Request)
	}
}
