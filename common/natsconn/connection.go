package natsconn

import (
	"handworks/common/utils"

	"github.com/nats-io/nats.go"
)

func ConnectNATS() *nats.Conn {
	logger, err := utils.NewLogger()
	if err != nil {
		panic(err)
	}
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		logger.Error("Error connecting to NATS: %s", err)
		return nil
	}
	return nc
}
