package service

import (
	"encoding/json"
	"fmt"
	"handworks/common/grpc/booking"
	types "handworks/common/types/booking"
	"time"

	"github.com/nats-io/nats.go"
)

func (b *BookingService) CollectResponses(sub *nats.Subscription) []types.BookingReply {
	var responses []types.BookingReply
	timeout := 2 * time.Second

	for {
		msg, err := sub.NextMsg(timeout)
		if err == nats.ErrTimeout {
			b.L.Info("No more replies received.")
			break
		}
		if err != nil {
			b.L.Error("Error receiving NATS reply: %v", err)
			break
		}

		var reply types.BookingReply
		if err := json.Unmarshal(msg.Data, &reply); err != nil {
			b.L.Info("Invalid JSON in reply: %v", err)
			continue
		}

		responses = append(responses, reply)
	}

	b.L.Info("Received %d replies", len(responses))
	return responses
}

func (b *BookingService) MergeBookingReplies(replies []types.BookingReply) (
	equipments []types.CleaningEquipment,
	resources []types.CleaningResources,
	cleaners []types.CleanerAssigned,
	prices types.CleaningPrices,
) {
	for _, r := range replies {
		equipments = append(equipments, r.Equipments...)
		resources = append(resources, r.Resources...)
		cleaners = append(cleaners, r.Cleaners...)
		prices.MainServicePrice = r.Prices.MainServicePrice
		prices.AddonPrices = append(prices.AddonPrices, r.Prices.AddonPrices...)
	}
	return
}
func (b *BookingService) ExtractAddonPrices(prices types.CleaningPrices) []types.AddonCleaningPrice {
	var addons []types.AddonCleaningPrice
	addons = append(addons, prices.AddonPrices...)
	return addons
}

func (b *BookingService) ExtractEquipmentIDs(equipments []*booking.CleaningEquipment) []string {
	var ids []string
	for _, e := range equipments {
		if e != nil && e.Id != "" {
			ids = append(ids, e.Id)
		}
	}
	return ids
}

func (b *BookingService) ExtractResourceIDs(resources []*booking.CleaningResources) []string {
	var ids []string
	for _, r := range resources {
		if r != nil && r.Id != "" {
			ids = append(ids, r.Id)
		}
	}
	return ids
}

func (b *BookingService) ExtractCleanerIDs(cleaners []*booking.CleanerAssigned) []string {
	var ids []string
	for _, c := range cleaners {
		if c != nil && c.Id != "" {
			ids = append(ids, c.Id)
		}
	}
	return ids
}

func (b *BookingService) createEquipment(
	// ctx context.Context,
	// tx pgx.Tx,
	equipment *booking.CleaningEquipment,
) (*types.CleaningEquipment, error) {
	var createdEquipment types.CleaningEquipment

	createdEquipment.ID = "d3b07384-4e6f-4f7e-bf6a-2e7f3c2c1f9a"
	createdEquipment.Name = equipment.Name
	createdEquipment.Type = equipment.Type
	createdEquipment.PhotoURL = equipment.PhotoUrl

	fmt.Print(createdEquipment)

	return &createdEquipment, nil
}

func (b *BookingService) createResource(
	// ctx context.Context,
	// tx pgx.Tx,
	resource *booking.CleaningResources,
) (*types.CleaningResources, error) {
	var createdResource types.CleaningResources

	createdResource.ID = "8f14e45f-ea9d-4c3d-9b1c-ff99a5b6d7c1"
	createdResource.Name = resource.Name
	createdResource.Type = resource.Type
	createdResource.PhotoURL = resource.PhotoUrl

	fmt.Print(createdResource)

	return &createdResource, nil
}

func (b *BookingService) createCleanersAssigned(
	// ctx context.Context,
	// tx pgx.Tx,
	cleaner *booking.CleanerAssigned,
) (*types.CleanerAssigned, error) {
	var createdCleaner types.CleanerAssigned

	createdCleaner.ID = "45c48cce-0f0e-4b8f-8c4e-9d1a7f2f5e2d"
	createdCleaner.CleanerFirstName = cleaner.CleanerFirstName
	createdCleaner.CleanerLastName = cleaner.CleanerLastName
	createdCleaner.PFPUrl = cleaner.PfpUrl

	fmt.Print(createdCleaner)

	return &createdCleaner, nil
}
