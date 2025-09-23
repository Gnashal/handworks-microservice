package types

import "handworks/common/grpc/booking"

type ServiceDetail struct {
	General  GeneralCleaningDetails
	Couch    CouchCleaningDetails
	Mattress MattressCleaningDetails
	Car      CarCleaningDetails
	Post     PostConstructionDetails
}

func ServiceDetailFromProto(pb *booking.ServiceDetail) ServiceDetail {
	if pb == nil {
		return ServiceDetail{}
	}

	sd := ServiceDetail{}

	switch t := pb.Type.(type) {
	case *booking.ServiceDetail_General:
		sd.General = GeneralCleaningDetailsFromProto(t.General)
	case *booking.ServiceDetail_Couch:
		sd.Couch = CouchCleaningDetailsFromProto(t.Couch)
	case *booking.ServiceDetail_Mattress:
		sd.Mattress = MattressCleaningDetailsFromProto(t.Mattress)
	case *booking.ServiceDetail_Car:
		sd.Car = CarCleaningDetailsFromProto(t.Car)
	case *booking.ServiceDetail_Post:
		sd.Post = PostConstructionCleaningDetailsFromProto(t.Post)
	}

	return sd
}

func (sd ServiceDetail) ToProto() *booking.ServiceDetail {
	pb := &booking.ServiceDetail{}

	if sd.General != (GeneralCleaningDetails{}) {
		pb.Type = &booking.ServiceDetail_General{
			General: sd.General.ToProto(),
		}
	} else if sd.Couch != (CouchCleaningDetails{}) {
		pb.Type = &booking.ServiceDetail_Couch{
			Couch: sd.Couch.ToProto(),
		}
	} else if sd.Mattress != (MattressCleaningDetails{}) {
		pb.Type = &booking.ServiceDetail_Mattress{
			Mattress: sd.Mattress.ToProto(),
		}
	} else if sd.Car != (CarCleaningDetails{}) {
		pb.Type = &booking.ServiceDetail_Car{
			Car: sd.Car.ToProto(),
		}
	} else if sd.Post != (PostConstructionDetails{}) {
		pb.Type = &booking.ServiceDetail_Post{
			Post: sd.Post.ToProto(),
		}
	}

	return pb
}

func (sd ServiceDetail) GetID() string {
	switch {
	case sd.General.ID != "":
		return sd.General.ID
	case sd.Couch.ID != "":
		return sd.Couch.ID
	case sd.Mattress.ID != "":
		return sd.Mattress.ID
	case sd.Car.ID != "":
		return sd.Car.ID
	case sd.Post.ID != "":
		return sd.Post.ID
	default:
		return ""
	}
}
