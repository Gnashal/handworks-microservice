package grpc

// add gprc config here
import (
	"fmt"
	accpb "handworks/common/grpc/account"
	bookpb "handworks/common/grpc/booking"
	intentorypb "handworks/common/grpc/inventory"
	paypb "handworks/common/grpc/payment"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClients struct {
	AccClient       accpb.AccountServiceClient
	BookingClient   bookpb.BookingServiceClient
	InventoryClient intentorypb.InventoryServiceClient
	PaymentClient   paypb.PaymentServiceClient
	// Add more gRPC clients here
}

func getServiceConnectionStrings() (string, string, string, string, error) {
	ACC_CONN := os.Getenv("ACC_CONN")
	BOOK_CONN := os.Getenv("BOOK_CONN")
	INV_CONN := os.Getenv("INV_CONN")
	PAY_CONN := os.Getenv("PAY_CONN")
	if ACC_CONN == "" {
		return "", "", "", "", fmt.Errorf("ACC_CONN environment variable not set")
	}
	if BOOK_CONN == "" {
		return "", "", "", "", fmt.Errorf("BOOK_CONN environment variable not set")
	}
	if INV_CONN == "" {
		return "", "", "", "", fmt.Errorf("INV_CONN environment variable not set")
	}
	if PAY_CONN == "" {
		return "", "", "", "", fmt.Errorf("PAY_CONN environment variable not set")
	}
	return ACC_CONN, BOOK_CONN, INV_CONN, PAY_CONN, nil
}

func NewClients() (*GrpcClients, error) {
	// Connection strings here
	ACC_CONN, BOOK_CONN, INV_CONN, PAY_CONN, err := getServiceConnectionStrings()
	if err != nil {
		return nil, err
	}
	// will changes this nya for prod, but for now this is fine
	accConn, err := grpc.NewClient(ACC_CONN,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	bookConn, err := grpc.NewClient(BOOK_CONN,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	invConn, err := grpc.NewClient(INV_CONN,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	payConn, err := grpc.NewClient(PAY_CONN,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &GrpcClients{
		AccClient:       accpb.NewAccountServiceClient(accConn),
		BookingClient:   bookpb.NewBookingServiceClient(bookConn),
		InventoryClient: intentorypb.NewInventoryServiceClient(invConn),
		PaymentClient:   paypb.NewPaymentServiceClient(payConn),
		// add more clients here
	}, nil
}
