package service

import (
	"context"
	"encoding/json"
	types "handworks/common/types/account"
	btypes "handworks/common/types/booking"

	"github.com/nats-io/nats.go"
)

func (i *AccountService) HandleSubscriptions(ctx context.Context) error {
	if err := i.SubscribeBookingRequests(); err != nil {
		i.L.Error("%v\n", err)
		return err
	}
	<-ctx.Done()
	return ctx.Err()

}

func (i *AccountService) SubscribeBookingRequests() error {
	_, err := i.NC.Subscribe("booking.created", func(msg *nats.Msg) {
		var req btypes.CreateBookingEvent
		if err := json.Unmarshal(msg.Data, &req); err != nil {
			i.L.Error("Failed to unmarshal booking request: %v\n", err)
			return
		}

		i.L.Info("Received booking.created event for customer: %s %s\n",
			req.Base.CustomerFirstName, req.Base.CustomerLastName)

		// FOR TESTING PA NI, I HAVE NOT IMPLEMENTED THE REAL LOGIC YET
		// TODO: Automation logic for cleaner assignment
		cleaners := []types.CleanerAssigned{
			{ID: "7aa794fa-e3f9-446f-8368-bcb55518bc29", CleanerFirstName: "Charles", CleanerLastName: "Boquecosa"},
			{ID: "cb32d23a-31a8-4461-ba3e-228d418ba6f3", CleanerFirstName: "Clarence", CleanerLastName: "Diangco"},
		}

		reply := types.BookingReply{
			Source:   "account-service",
			Cleaners: cleaners,
		}

		data, err := json.Marshal(reply)
		if err != nil {
			i.L.Error("Failed to marshal reply: %v\n", err)
			return
		}

		if msg.Reply != "" {
			if err := i.NC.Publish(msg.Reply, data); err != nil {
				i.L.Error("Failed to publish reply: %v\n", err)
			} else {
				i.L.Info("Sent booking reply from account-service")
			}
		} else {
			i.L.Info("No reply subject — skipping reply")
		}
	})
	if err != nil {
		return err
	}

	i.L.Info("Account service subscribed to booking.created")
	return nil
}
