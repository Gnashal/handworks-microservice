package config

import (
	"context"
	graph "handworks-gateway/graph/generated"
	"handworks-gateway/graph/resolvers"
	"handworks-gateway/grpc"
	"handworks-gateway/internal/middleware"
	"handworks/common/utils"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
)

func StartGQlServer(l *utils.Logger, wg *sync.WaitGroup, stop <-chan struct{}) {
	defer wg.Done()
	clerkKey := os.Getenv("CLERK_SECRET_KEY")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	clients, err := grpc.NewClients()
	if err != nil {
		l.Fatal("Failed to dial grpc servers: %w", err)
	}
	config := graph.Config{
		Resolvers: &resolvers.Resolver{
			GrpcClients: clients,
		},
	}
	srv := handler.New(graph.NewExecutableSchema(config))
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	cache := lru.New[*ast.QueryDocument](1000)
	srv.SetQueryCache(cache)

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: setupRoutes(srv, clerkKey),
	}

	go func() {
		<-stop
		l.Info("Shutting down GraphQL server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = httpServer.Shutdown(ctx)
	}()
	l.Info("GrahQL playground running on http://localhost:%s/playground", port)
	l.Info("GraphQL server running on http://localhost:%s/api", port)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		l.Error("GraphQL server error: %v", err)
	}
}

// helper to register routes
func setupRoutes(srv *handler.Server, clerkKey string) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/playground", playground.Handler("GraphQL playground", "/api"))
	mux.Handle("/api", middleware.AuthMiddleware(srv, clerkKey))
	return mux
}
