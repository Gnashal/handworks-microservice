package services

import (
	"context"
	"encoding/json"
	"fmt"
	"handworks-api/tasks"
	"handworks-api/types"
	"time"

	"github.com/google/uuid"
)

// func (s *BookingService) withTx(
//
//	ctx context.Context,
//	fn func(pgx.Tx) error,
//
//	) (err error) {
//		tx, err := s.DB.Begin(ctx)
//		if err != nil {
//			return fmt.Errorf("failed to begin tx: %w", err)
//		}
//		defer func() {
//			if err != nil {
//				if rbErr := tx.Rollback(ctx); rbErr != nil {
//					s.Logger.Error("rollback failed: %v", rbErr)
//				}
//			} else {
//				err = tx.Commit(ctx)
//			}
//		}()
//		return fn(tx)
//	}
func (s *BookingService) CreateBooking(ctx context.Context, evt types.CreateBookingEvent) (*types.Booking, error) {
	if evt.Base.CustID == "" {
		return nil, fmt.Errorf("missing customer id")
	}

	now := time.Now()
	if evt.Base.ID == "" {
		evt.Base.ID = uuid.NewString()
	}
	evt.Base.CreatedAt = now
	evt.Base.UpdatedAt = nil

	if evt.Base.Schedule.ID == "" {
		evt.Base.Schedule.ID = uuid.NewString()
	}

	materialize := func(req types.ServicesRequest) (types.ServiceDetails, error) {
		sd := types.ServiceDetails{
			ID:          uuid.NewString(),
			ServiceType: string(req.ServiceType),
		}

		dt := types.DetailType(req.ServiceType)
		if f, ok := types.DetailFactories[dt]; ok {
			target := f()
			b, err := json.Marshal(req.Details)
			if err != nil {
				return sd, err
			}
			if err := json.Unmarshal(b, target); err != nil {
				return sd, err
			}
			sd.Details = target
			return sd, nil
		}

		sd.Details = req.Details
		return sd, nil
	}

	mainSvc, err := materialize(evt.MainService)
	if err != nil {
		return nil, fmt.Errorf("materialize main service: %w", err)
	}

	buildEquipmentsForService := func(sd types.ServiceDetails) []types.CleaningEquipment {
		var baseEquipment []types.CleaningEquipment
		switch types.DetailType(sd.ServiceType) {
		case types.ServiceGeneral:
			if d, ok := sd.Details.(*types.GeneralCleaningDetails); ok {
				baseEquipment = append(baseEquipment, types.CleaningEquipment{ID: uuid.NewString(), Name: "Vacuum Cleaner", Type: "MACHINE"})
				baseEquipment = append(baseEquipment, types.CleaningEquipment{ID: uuid.NewString(), Name: "Mop & Bucket", Type: "TOOL"})
				if d.SQM > 50 {
					baseEquipment = append(baseEquipment, types.CleaningEquipment{ID: uuid.NewString(), Name: "Extra Mop Set", Type: "TOOL"})
				}
			}
		case types.ServiceCouch:
			if d, ok := sd.Details.(*types.CouchCleaningDetails); ok {
				for _, spec := range d.CleaningSpecs {
					for i := int32(0); i < spec.Quantity; i++ {
						baseEquipment = append(baseEquipment, types.CleaningEquipment{ID: uuid.NewString(), Name: "Upholstery Cleaner", Type: "MACHINE"})
					}
				}
				if d.BedPillows > 0 {
					baseEquipment = append(baseEquipment, types.CleaningEquipment{ID: uuid.NewString(), Name: "Pillow Vacuum", Type: "TOOL"})
				}
			}
		case types.ServiceMattress:
			if d, ok := sd.Details.(*types.MattressCleaningDetails); ok {
				for _, spec := range d.CleaningSpecs {
					for i := int32(0); i < spec.Quantity; i++ {
						baseEquipment = append(baseEquipment, types.CleaningEquipment{ID: uuid.NewString(), Name: "Mattress Steam Cleaner", Type: "MACHINE"})
					}
				}
			}
		case types.ServiceCar:
			if d, ok := sd.Details.(*types.CarCleaningDetails); ok {
				for _, spec := range d.CleaningSpecs {
					for i := int32(0); i < spec.Quantity; i++ {
						baseEquipment = append(baseEquipment, types.CleaningEquipment{ID: uuid.NewString(), Name: "Car Shampoo Machine", Type: "MACHINE"})
					}
				}
				if d.ChildSeats > 0 {
					baseEquipment = append(baseEquipment, types.CleaningEquipment{ID: uuid.NewString(), Name: "Child Seat Cleaner", Type: "TOOL"})
				}
			}
		case types.ServicePost:
			if d, ok := sd.Details.(*types.PostConstructionDetails); ok {
				baseEquipment = append(baseEquipment, types.CleaningEquipment{ID: uuid.NewString(), Name: "Industrial Vacuum", Type: "MACHINE"})
				if d.SQM > 100 {
					baseEquipment = append(baseEquipment, types.CleaningEquipment{ID: uuid.NewString(), Name: "Extra Industrial Vacuum", Type: "MACHINE"})
				}
			}
		default:
			baseEquipment = append(baseEquipment, types.CleaningEquipment{ID: uuid.NewString(), Name: "Standard Cleaning Kit", Type: "KIT"})
		}
		return baseEquipment
	}

	buildResourcesForService := func(sd types.ServiceDetails) []types.CleaningResources {
		var res []types.CleaningResources
		switch types.DetailType(sd.ServiceType) {
		case types.ServiceGeneral:
			if d, ok := sd.Details.(*types.GeneralCleaningDetails); ok {
				res = append(res, types.CleaningResources{ID: uuid.NewString(), Name: "All-purpose Cleaner", Type: "LIQUID"})
				res = append(res, types.CleaningResources{ID: uuid.NewString(), Name: "Garbage Bags", Type: "SUPPLY"})
				if d.SQM > 50 {
					res = append(res, types.CleaningResources{ID: uuid.NewString(), Name: "Extra Cleaning Solution", Type: "LIQUID"})
				}
			}
		case types.ServiceCouch:
			if d, ok := sd.Details.(*types.CouchCleaningDetails); ok {
				res = append(res, types.CleaningResources{ID: uuid.NewString(), Name: "Upholstery Shampoo", Type: "LIQUID"})
				if d.BedPillows > 0 {
					res = append(res, types.CleaningResources{ID: uuid.NewString(), Name: "Pillow Covers", Type: "SUPPLY"})
				}
			}
		case types.ServiceMattress:
			res = append(res, types.CleaningResources{ID: uuid.NewString(), Name: "Mattress Cleaner", Type: "LIQUID"})
		case types.ServiceCar:
			res = append(res, types.CleaningResources{ID: uuid.NewString(), Name: "Car Shampoo", Type: "LIQUID"})
		case types.ServicePost:
			res = append(res, types.CleaningResources{ID: uuid.NewString(), Name: "Debris Bags", Type: "SUPPLY"})
		default:
			res = append(res, types.CleaningResources{ID: uuid.NewString(), Name: "Standard Supplies", Type: "SUPPLY"})
		}
		return res
	}
	// --- end helpers ---

	priceCalc := &tasks.PaymentTasks{}
	mainPrice := priceCalc.CalculatePriceByServiceType(&evt.MainService)

	// populate equipments/resources for main service
	equipments := buildEquipmentsForService(mainSvc)
	resources := buildResourcesForService(mainSvc)

	// addons
	var addons []types.AddOns
	var addonTotal float32 = 0
	for _, a := range evt.Addons {
		sdet, err := materialize(a.ServiceDetail)
		if err != nil {
			return nil, fmt.Errorf("materialize addon: %w", err)
		}

		price := priceCalc.CalculatePriceByServiceType(&a.ServiceDetail)
		addonTotal += price

		// merge addon-specific equipments/resources (append to main lists)
		addEquip := buildEquipmentsForService(sdet)
		addRes := buildResourcesForService(sdet)
		equipments = append(equipments, addEquip...)
		resources = append(resources, addRes...)

		addons = append(addons, types.AddOns{
			ID:            uuid.NewString(),
			ServiceDetail: sdet,
			Price:         price,
		})
	}

	booking := &types.Booking{
		ID:          uuid.NewString(),
		Base:        evt.Base,
		MainService: mainSvc,
		Schedule:    evt.Base.Schedule,
		Addons:      addons,
		Equipments:  equipments,
		Resources:   resources,
		Cleaners:    nil, //TODO: CLEANER ASSIGNMENT LOGIC
		TotalPrice:  mainPrice + addonTotal,
	}
	return booking, nil
}

func (s *BookingService) GetBookingById(ctx context.Context) error {
	return nil
}

func (s *BookingService) GetBookingByUId(ctx context.Context) error {
	return nil
}

func (s *BookingService) UpdateBooking(ctx context.Context) error {
	return nil
}

func (s *BookingService) DeleteBooking(ctx context.Context) error {
	return nil
}
