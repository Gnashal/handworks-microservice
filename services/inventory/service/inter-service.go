package service

import (
	"context"
	"encoding/json"
	btypes "handworks/common/types/booking"
	types "handworks/common/types/inventory"

	"github.com/nats-io/nats.go"
)

func (i *InventoryService) HandleSubscriptions(ctx context.Context) error {
	if err := i.SubscribeBookingRequests(); err != nil {
		i.L.Error("%v\n", err)
		return err
	}
	<-ctx.Done()
	return ctx.Err()

}

func (i *InventoryService) SubscribeBookingRequests() error {
	_, err := i.NC.Subscribe("booking.created", func(msg *nats.Msg) {
		var req btypes.CreateBookingEvent
		if err := json.Unmarshal(msg.Data, &req); err != nil {
			i.L.Error("Failed to unmarshal booking request: %v\n", err)
			return
		}

		i.L.Info("Received booking.created event for customer: %s %s\n",
			req.Base.CustomerFirstName, req.Base.CustomerLastName)

		// FOR TESTING PA NI, I HAVE NOT IMPLEMENTED THE REAL LOGIC YET
		// TODO: Automation logic for resource and equipment allocation
		equipments := []types.CleaningEquipment{
			{ID: "7849f478-f70b-42a7-82d2-aadc81d3e6d6", Name: "Vacuum Cleaner", Type: "Electrical", PhotoURL: "https://example.com/vacuum.jpg"},
			{ID: "a4cd7e23-787b-4344-80d7-c50199d85ecd", Name: "Mop", Type: "Manual", PhotoURL: "https://example.com/mop.jpg"},
		}
		resources := []types.CleaningResources{
			{ID: "d1e94940-838d-4916-bf2b-bb09b77d7c46", Name: "Detergent", Type: "Chemical", PhotoURL: "https://example.com/detergent.jpg"},
		}

		reply := types.BookingReply{
			Source:     "inventory-service",
			Equipments: equipments,
			Resources:  resources,
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
				i.L.Info("Sent booking reply from inventory-service")
			}
		} else {
			i.L.Info("No reply subject — skipping reply")
		}
	})
	if err != nil {
		return err
	}

	i.L.Info("Inventory service subscribed to booking.created")
	return nil
}
