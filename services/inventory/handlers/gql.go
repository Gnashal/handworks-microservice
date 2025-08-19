package handlers

import (
	"context"
	"handworks/common/utils"
	"handworks/services/inventory/graph"
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

// grpc handlers for the  admin service
func StartGQlServer(l *utils.Logger, wg *sync.WaitGroup, stop <-chan struct{}) {
	defer wg.Done()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
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
		Handler: setupRoutes(srv),
	}

	go func() {
		<-stop
		l.Info("Shutting down GraphQL server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = httpServer.Shutdown(ctx)
	}()
	l.Info("GrahQL playground running on http://localhost:%s/playground", port)
	l.Info("GraphQL server running on http://localhost:%s", port)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		l.Error("GraphQL server error: %v", err)
	}
}

// helper to register routes
func setupRoutes(srv *handler.Server) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/playground", playground.Handler("GraphQL playground", "/query"))
	mux.Handle("/query", srv)
	return mux
}
