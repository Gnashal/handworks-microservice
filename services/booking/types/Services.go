package types

import (
	"fmt"
	"handworks/common/grpc/booking"
)

type ServiceDetails struct {
	ID          string
	ServiceType string
	Details     any
}

func (mainService ServiceDetails) ToProto() *booking.Services {
	var detail *booking.ServiceDetail

	switch types := mainService.Details.(type) {
	case *GeneralCleaningDetails:
		detail = &booking.ServiceDetail{Type: &booking.ServiceDetail_General{General: types.ToProto()}}
	case *CouchCleaningDetails:
		detail = &booking.ServiceDetail{Type: &booking.ServiceDetail_Couch{Couch: types.ToProto()}}
	case *MattressCleaningDetails:
		detail = &booking.ServiceDetail{Type: &booking.ServiceDetail_Mattress{Mattress: types.ToProto()}}
	case *CarCleaningDetails:
		detail = &booking.ServiceDetail{Type: &booking.ServiceDetail_Car{Car: types.ToProto()}}
	case *PostConstructionDetails:
		detail = &booking.ServiceDetail{Type: &booking.ServiceDetail_Post{Post: types.ToProto()}}
	default:
		detail = nil
	}

	return &booking.Services{
		Id:          mainService.ID,
		ServiceType: booking.MainServiceType(booking.MainServiceType_value[mainService.ServiceType]),
		Details:     detail,
	}
}

func ServiceFromProto(pb *booking.Services) ServiceDetails {
	if pb == nil {
		return ServiceDetails{}
	}

	var details any

	switch d := pb.Details.Type.(type) {
	case *booking.ServiceDetail_General:
		details = GeneralCleaningDetailsFromProto(d.General)
	case *booking.ServiceDetail_Couch:
		details = CouchCleaningDetailsFromProto(d.Couch)
	case *booking.ServiceDetail_Mattress:
		details = MattressCleaningDetailsFromProto(d.Mattress)
	case *booking.ServiceDetail_Car:
		details = CarCleaningDetailsFromProto(d.Car)
	case *booking.ServiceDetail_Post:
		details = PostConstructionCleaningDetailsFromProto(d.Post)
	}

	return ServiceDetails{
		ID:          pb.Id,
		ServiceType: pb.ServiceType.String(),
		Details:     details,
	}
}

func UnmarshalServiceDetails(serviceType string, raw []byte) (any, error) {
	switch serviceType {
	case "GENERAL":
		return UnmarshalGeneralDetails(raw)
	case "COUCH":
		return UnmarshalCouchDetails(raw)
	case "MATTRESS":
		return UnmarshalMattressDetails(raw)
	case "CAR":
		return UnmarshalCarDetails(raw)
	case "POST":
		return UnmarshalPostConstructionDetails(raw)
	default:
		return nil, fmt.Errorf("unsupported service type: %s", serviceType)
	}
}
