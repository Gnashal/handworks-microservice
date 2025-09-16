package grpc

// add gprc config here
import (
	accpb "handworks/common/grpc/account"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClients struct {
	AccClient accpb.AccountServiceClient
	// Add more gRPC clients here
}

func NewClients() (*GrpcClients, error) {
	// Connection strings here
	ACC_CONN := os.Getenv("ACC_CONN")

	// Then add the connections here
	accConn, err := grpc.NewClient(ACC_CONN,
		// will changes this nya for prod, but for now this is fine
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &GrpcClients{
		AccClient: accpb.NewAccountServiceClient(accConn),
		// add more clients here
	}, nil
}
