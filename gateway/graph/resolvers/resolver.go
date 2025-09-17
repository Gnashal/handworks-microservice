package resolvers

//go:generate go run github.com/99designs/gqlgen generate
import (
	"handworks-gateway/grpc"
	"handworks/common/utils"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Log         *utils.Logger
	GrpcClients *grpc.GrpcClients
}
