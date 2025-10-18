package helpers

import (
	"fmt"
	model "handworks-gateway/graph/generated/models"
	"handworks/common/grpc/booking"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func MapServiceDetail(in *booking.ServiceDetail) *model.ServiceDetail {
	if in == nil {
		return nil
	}

	sd := &model.ServiceDetail{}

	switch detail := in.Type.(type) {
	case *booking.ServiceDetail_General:
		sd.General = &model.GeneralCleaningDetails{
			HomeType: detail.General.HomeType.String(),
			Sqm:      int32(detail.General.Sqm),
		}
	case *booking.ServiceDetail_Couch:
		var specs []*model.CouchCleaningSpecifications
		for _, s := range detail.Couch.CleaningSpecs {
			specs = append(specs, &model.CouchCleaningSpecifications{
				CouchType: s.CouchType.String(),
				WidthCm:   int32(s.WidthCm),
				DepthCm:   int32(s.DepthCm),
				HeightCm:  int32(s.HeightCm),
				Quantity:  int32(s.Quantity),
			})
		}
		sd.Couch = &model.CouchCleaningDetails{
			CleaningSpecs: specs,
			BedPillows:    ptrInt(int32(detail.Couch.BedPillows)),
		}
	case *booking.ServiceDetail_Mattress:
		var specs []*model.MattressCleaningSpecifications
		for _, s := range detail.Mattress.CleaningSpecs {
			specs = append(specs, &model.MattressCleaningSpecifications{
				BedType:  s.BedType.String(),
				WidthCm:  int32(s.WidthCm),
				DepthCm:  int32(s.DepthCm),
				HeightCm: int32(s.HeightCm),
				Quantity: int32(s.Quantity),
			})
		}
		sd.Mattress = &model.MattressCleaningDetails{
			CleaningSpecs: specs,
		}
	case *booking.ServiceDetail_Car:
		var specs []*model.CarCleaningSpecifications
		for _, s := range detail.Car.CleaningSpecs {
			specs = append(specs, &model.CarCleaningSpecifications{
				CarType:  s.CarType.String(),
				Quantity: int32(s.Quantity),
			})
		}
		sd.Car = &model.CarCleaningDetails{
			CleaningSpecs: specs,
			ChildSeats:    ptrInt(int32(detail.Car.ChildSeats)),
		}
	case *booking.ServiceDetail_Post:
		sd.Post = &model.PostConstructionCleaningDetails{
			Sqm: int32(detail.Post.Sqm),
		}
	}

	return sd
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
				BedPillows:    int32(input.Details.Couch.BedPillows),
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
func MapAddressModelToProto(in *model.AddressInput) *booking.Address {
	if in == nil {
		return nil
	}
	return &booking.Address{
		AddressHuman: in.AddressHuman,
		AddressLat:   in.AddressLat,
		AddressLng:   in.AddressLng,
	}
}
func MapAddresTypeModelToProto(add *model.Address) *booking.Address {
	return &booking.Address{
		AddressHuman: add.AddressHuman,
		AddressLat:   add.AddressLat,
		AddressLng:   add.AddressLng,
	}
}

func MapAddressProtoToModel(in *booking.Address) *model.Address {
	if in == nil {
		return nil
	}
	return &model.Address{
		AddressHuman: in.AddressHuman,
		AddressLat:   in.AddressLat,
		AddressLng:   in.AddressLng,
	}
}

func MapAddOnInput(input *model.AddOnInput) *booking.AddOnRequest {
	if input == nil || input.Service == nil {
		return nil
	}

	return &booking.AddOnRequest{
		ServiceDetail: MapServicesInput(input.Service),
	}
}
func MapAddOnInputs(inputs []*model.AddOnInput) []*booking.AddOnRequest {
	if len(inputs) == 0 {
		return nil
	}

	addons := make([]*booking.AddOnRequest, 0, len(inputs))
	for _, a := range inputs {
		mapped := MapAddOnInput(a)
		if mapped != nil {
			addons = append(addons, mapped)
		}
	}
	return addons
}

func MapBaseDetailsInputToProto(in *model.BaseBookingDetailsInput) *booking.BaseBookingDetailsRequest {
	if in == nil {
		return nil
	}
	return &booking.BaseBookingDetailsRequest{
		CustId:            in.CustID,
		CustomerFirstName: in.CustomerFirstName,
		CustomerLastName:  in.CustomerLastName,
		Address:           MapAddressModelToProto(in.Address),
		StartSched:        timestamppb.New(in.StartSched),
		EndSched:          timestamppb.New(in.EndSched),
		DirtyScale:        int32(in.DirtyScale),
		PaymentStatus:     in.PaymentStatus,
		ReviewStatus:      in.ReviewStatus,
		Photos:            in.Photos,
		QuoteId:           *in.QuoteID,
	}
}
func MapBaseDetailsTypeToProto(in *model.BaseBookingDetails) *booking.BaseBookingDetails {
	if in == nil {
		return nil
	}
	return &booking.BaseBookingDetails{
		CustId:            in.CustID,
		CustomerFirstName: in.CustomerFirstName,
		CustomerLastName:  in.CustomerLastName,
		Address:           MapAddresTypeModelToProto(in.Address),
		StartSched:        timestamppb.New(in.StartSched),
		EndSched:          timestamppb.New(in.EndSched),
		DirtyScale:        int32(in.DirtyScale),
		PaymentStatus:     in.PaymentStatus,
		ReviewStatus:      in.ReviewStatus,
		Photos:            in.Photos,
		QuoteId:           *in.QuoteID,
	}
}

func MapBaseDetailsProtoToModel(in *booking.BaseBookingDetails) *model.BaseBookingDetails {
	if in == nil {
		return nil
	}
	return &model.BaseBookingDetails{
		ID:                in.Id,
		CustID:            in.CustId,
		CustomerFirstName: in.CustomerFirstName,
		CustomerLastName:  in.CustomerLastName,
		Address:           MapAddressProtoToModel(in.Address),
		StartSched:        in.StartSched.AsTime(),
		EndSched:          in.EndSched.AsTime(),
		DirtyScale:        int32(in.DirtyScale),
		PaymentStatus:     in.PaymentStatus,
		ReviewStatus:      in.ReviewStatus,
		Photos:            in.Photos,
		QuoteID:           &in.QuoteId,
	}
}
func MapServicesInput(input *model.ServicesInput) *booking.ServicesRequest {
	if input == nil {
		return nil
	}

	return &booking.ServicesRequest{
		ServiceType: MapMainServiceType(input.ServiceType),
		Details:     MapServiceDetailInput(input),
	}
}
func MapServices(in *booking.Services) *model.Services {
	if in == nil {
		return nil
	}

	return &model.Services{
		ID:          in.Id,
		ServiceType: in.ServiceType.String(),
		Details:     MapServiceDetail(in.Details),
	}
}

func MapAddOn(in *booking.AddOn) *model.AddOn {
	if in == nil {
		return nil
	}

	return &model.AddOn{
		ID:            in.Id,
		Price:         float64(in.Price),
		ServiceDetail: MapServices(in.ServiceDetail),
	}
}

func MapBooking(in *booking.Booking) *model.Booking {
	if in == nil {
		fmt.Print("Addons is nil when returned")
		return nil
	}

	return &model.Booking{
		ID:          in.Id,
		Base:        MapBaseDetailsProtoToModel(in.Base),
		MainService: MapServices(in.MainService),
		Addons:      MapAddOns(in.Addons),
		Equipment:   MapCleaningEquipmentList(in.Equipment),
		Resources:   MapCleaningResourcesList(in.Resources),
		Cleaners:    MapCleanerAssignedList(in.Cleaners),
		TotalPrice:  float64(in.TotalPrice),
	}
}
func MapAddOns(in []*booking.AddOn) []*model.AddOn {
	if len(in) == 0 {
		return nil
	}
	addons := make([]*model.AddOn, 0, len(in))
	for _, a := range in {
		addons = append(addons, MapAddOn(a))
	}
	return addons
}
func MapCleaningEquipment(in *booking.CleaningEquipment) *model.CleaningEquipment {
	return &model.CleaningEquipment{
		ID:       in.Id,
		Name:     in.Name,
		Type:     in.Type,
		PhotoURL: &in.PhotoUrl,
	}
}
func MapCleaningResources(in *booking.CleaningResources) *model.CleaningResources {
	return &model.CleaningResources{
		ID:       in.Id,
		Name:     in.Name,
		Type:     in.Type,
		PhotoURL: &in.PhotoUrl,
	}
}
func MapCleanerAssigned(in *booking.CleanerAssigned) *model.CleanerAssigned {
	return &model.CleanerAssigned{
		ID:               in.Id,
		CleanerFirstName: in.CleanerFirstName,
		CleanerLastName:  in.CleanerLastName,
		PfpURL:           &in.PfpUrl,
	}
}

func MapCleaningEquipmentList(in []*booking.CleaningEquipment) []*model.CleaningEquipment {
	if len(in) == 0 {
		return nil
	}
	equipment := make([]*model.CleaningEquipment, 0, len(in))
	for _, e := range in {
		equipment = append(equipment, MapCleaningEquipment(e))
	}
	return equipment
}

func MapCleaningResourcesList(in []*booking.CleaningResources) []*model.CleaningResources {
	if len(in) == 0 {
		return nil
	}
	resources := make([]*model.CleaningResources, 0, len(in))
	for _, r := range in {
		resources = append(resources, MapCleaningResources(r))
	}
	return resources
}

func MapCleanerAssignedList(in []*booking.CleanerAssigned) []*model.CleanerAssigned {
	if len(in) == 0 {
		return nil
	}
	cleaners := make([]*model.CleanerAssigned, 0, len(in))
	for _, c := range in {
		cleaners = append(cleaners, MapCleanerAssigned(c))
	}
	return cleaners
}
