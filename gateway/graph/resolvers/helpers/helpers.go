package helpers

import (
	model "handworks-gateway/graph/generated/models"
	"handworks/common/grpc/account"
	"handworks/common/grpc/booking"
	"handworks/common/grpc/inventory"
	"handworks/common/grpc/payment"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// Account helpers
func MapAccount(a *account.Account) *model.Account {
	if a == nil {
		return nil
	}

	return &model.Account{
		ID:        a.Id,
		FirstName: a.FirstName,
		LastName:  a.LastName,
		Email:     a.Email,
		Provider:  &a.Provider,
		Role:      a.Role,
		ClerkID:   a.ClerkId,
		CreatedAt: timestampToTime(a.CreatedAt),
		UpdatedAt: timestampToTime(a.UpdatedAt),
	}
}

func MapCustomer(c *account.Customer) *model.Customer {
	if c == nil {
		return nil
	}

	return &model.Customer{
		ID:      c.Id,
		Account: MapAccount(c.Account),
	}
}

func MapEmployee(e *account.Employee) *model.Employee {
	if e == nil {
		return nil
	}

	return &model.Employee{
		ID:               e.Id,
		Account:          MapAccount(e.Account),
		Position:         e.Position,
		Status:           e.Status.String(),
		PerformanceScore: float64(e.PerformanceScore),
		HireDate:         timestampToTime(e.HireDate),
		NumRatings:       int32(e.NumRatings),
	}
}

// func MapSignUpCustomer(e* acc)

// Convert protobuf Timestamp to Go time.Time
func timestampToTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}
	return ts.AsTime()
}

// Payment helpers

func MapQuote(res *payment.QuoteResponse) *model.Quote {
	if res == nil {
		return nil
	}

	addons := make([]*model.AddOnBreakdown, len(res.AddonBreakdown))
	for i, a := range res.AddonBreakdown {
		addons[i] = &model.AddOnBreakdown{
			AddonID:   a.AddonId,
			AddonName: a.AddonName,
			Price:     float64(a.Price),
		}
	}

	return &model.Quote{
		QuoteID:          res.QuoteId,
		MainServiceName:  res.MainServiceName,
		MainServiceTotal: float64(res.MainServiceTotal),
		AddonBreakdown:   addons,
		AddonTotal:       float64(res.AddonTotal),
		TotalPrice:       float64(res.TotalPrice),
	}
}

func MapServiceDetailInput(input *model.ServicesInput) *booking.ServiceDetail {
	if input == nil || input.Details == nil {
		return nil
	}

	sd := &booking.ServiceDetail{}

	switch input.ServiceType {
	case "GENERAL_CLEANING":
		sd.Type = &booking.ServiceDetail_General{
			General: &booking.GeneralCleaningDetails{
				HomeType: MapHomeType(input.Details.General.HomeType),
				Sqm:      int32(input.Details.General.Sqm),
			},
		}
	case "COUCH":
		var specs []*booking.CouchCleaningSpecifications
		for _, s := range input.Details.Couch.CleaningSpecs {
			specs = append(specs, &booking.CouchCleaningSpecifications{
				CouchType: MapCouchType(s.CouchType),
				WidthCm:   int32(s.WidthCm),
				DepthCm:   int32(s.DepthCm),
				HeightCm:  int32(s.HeightCm),
				Quantity:  int32(s.Quantity),
			})
		}
		sd.Type = &booking.ServiceDetail_Couch{
			Couch: &booking.CouchCleaningDetails{
				CleaningSpecs: specs,
				BedPillows:    int32(*input.Details.Couch.BedPillows),
			},
		}
	case "MATTRESS":
		var specs []*booking.MattressCleaningSpecifications
		for _, s := range input.Details.Mattress.CleaningSpecs {
			specs = append(specs, &booking.MattressCleaningSpecifications{
				BedType:  MapBedType(s.BedType),
				WidthCm:  int32(s.WidthCm),
				DepthCm:  int32(s.DepthCm),
				HeightCm: int32(s.HeightCm),
				Quantity: int32(s.Quantity),
			})
		}
		sd.Type = &booking.ServiceDetail_Mattress{
			Mattress: &booking.MattressCleaningDetails{
				CleaningSpecs: specs,
			},
		}
	case "CAR":
		var specs []*booking.CarCleaningSpecifications
		for _, s := range input.Details.Car.CleaningSpecs {
			specs = append(specs, &booking.CarCleaningSpecifications{
				CarType:  MapCarType(s.CarType),
				Quantity: int32(s.Quantity),
			})
		}
		sd.Type = &booking.ServiceDetail_Car{
			Car: &booking.CarCleaningDetails{
				CleaningSpecs: specs,
				ChildSeats:    int32(*input.Details.Car.ChildSeats),
			},
		}
	case "POST":
		sd.Type = &booking.ServiceDetail_Post{
			Post: &booking.PostConstructionCleaningDetails{
				Sqm: int32(input.Details.Post.Sqm),
			},
		}
	}

	return sd
}

// Booking helpers
func MapMainServiceType(s string) booking.MainServiceType {
	switch s {
	case "GENERAL_CLEANING":
		return booking.MainServiceType_GENERAL_CLEANING
	case "COUCH":
		return booking.MainServiceType_COUCH
	case "MATTRESS":
		return booking.MainServiceType_MATTRESS
	case "CAR":
		return booking.MainServiceType_CAR
	case "POST":
		return booking.MainServiceType_POST
	default:
		return booking.MainServiceType_SERVICE_TYPE_UNSPECIFIED
	}
}

func MapHomeType(s string) booking.HomeType {
	switch s {
	case "CONDO_ROOM":
		return booking.HomeType_CONDO_ROOM
	case "HOUSE":
		return booking.HomeType_HOUSE
	case "COMMERCIAL":
		return booking.HomeType_COMMERCIAL
	default:
		return booking.HomeType(0)
	}
}

func MapCouchType(s string) booking.CouchType {
	switch s {
	case "SEATER_1":
		return booking.CouchType_SEATER_1
	case "SEATER_2":
		return booking.CouchType_SEATER_2
	case "SEATER_3":
		return booking.CouchType_SEATER_3
	case "SEATER_3_LTYPE_SMALL":
		return booking.CouchType_SEATER_3_LTYPE_SMALL
	case "SEATER_3_LTYPE_LARGE":
		return booking.CouchType_SEATER_3_LTYPE_LARGE
	case "SEATER_4_LTYPE_SMALL":
		return booking.CouchType_SEATER_4_LTYPE_SMALL
	case "SEATER_4_LTYPE_LARGE":
		return booking.CouchType_SEATER_4_LTYPE_LARGE
	case "SEATER_5_LTYPE":
		return booking.CouchType_SEATER_5_LTYPE
	case "SEATER_6_LTYPE":
		return booking.CouchType_SEATER_6_LTYPE
	case "OTTOMAN":
		return booking.CouchType_OTTOMAN
	case "LAZBOY":
		return booking.CouchType_LAZBOY
	case "CHAIR":
		return booking.CouchType_CHAIR
	default:
		return booking.CouchType(0)
	}
}

func MapBedType(s string) booking.BedType {
	switch s {
	case "KING":
		return booking.BedType_KING
	case "QUEEN":
		return booking.BedType_QUEEN
	case "KING_HEADBAND":
		return booking.BedType_KING_HEADBAND
	case "QUEEN_HEADBAND":
		return booking.BedType_QUEEN_HEADBAND
	case "DOUBLE":
		return booking.BedType_DOUBLE
	case "SINGLE":
		return booking.BedType_SINGLE
	case "BED_PILLOW":
		return booking.BedType_BED_PILLOW
	default:
		return booking.BedType(0)
	}
}

func MapCarType(s string) booking.CarType {
	switch s {
	case "SEDAN":
		return booking.CarType_SEDAN
	case "MPV":
		return booking.CarType_MPV
	case "SUV":
		return booking.CarType_SUV
	case "PICKUP":
		return booking.CarType_PICKUP
	case "VAN":
		return booking.CarType_VAN
	case "CAR_SMALL":
		return booking.CarType_CAR_SMALL
	default:
		return booking.CarType(0)
	}
}

// Inventory helpers

func MapInventoryItem(in *inventory.InventoryItem) *model.InventoryItem {
	if in == nil {
		return nil
	}
	return &model.InventoryItem{
		ID:          in.Id,
		Name:        in.Name,
		Type:        in.Type.String(),
		Status:      in.Status.String(),
		Category:    in.Status.String(),
		Quantity:    in.Quantity,
		MaxQuantity: &in.MaxQuantity,
		Unit:        &in.Unit,
		IsAvailable: in.IsAvailable,
		ImageURL:    &in.ImageUrl,
		CreatedAt:   timestampToTime(in.CreatedAt),
		UpdatedAt:   timestampToTime(in.UpdatedAt),
	}
}

func MapInventoryItems(itemsIn []*inventory.InventoryItem) []*model.InventoryItem {

	items := make([]*model.InventoryItem, 0, len(itemsIn))
	for _, item := range itemsIn {
		items = append(items, MapInventoryItem(item))
	}
	return items
}
