package types

import (
	"encoding/json"
	"fmt"
	"handworks/common/grpc/booking"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (generalCleaning GeneralCleaningDetails) ToProto() *booking.GeneralCleaningDetails {
	return &booking.GeneralCleaningDetails{
		HomeType: booking.HomeType(booking.HomeType_value[generalCleaning.HomeType]),
		Sqm:      generalCleaning.SQM,
	}
}

func GeneralCleaningDetailsFromProto(pb *booking.GeneralCleaningDetails) *GeneralCleaningDetails {
	if pb == nil {
		return &GeneralCleaningDetails{}
	}
	return &GeneralCleaningDetails{
		HomeType: pb.HomeType.String(),
		SQM:      pb.Sqm,
	}
}

func (generalCleaning *GeneralCleaningDetails) MarshalGeneralDetails() ([]byte, error) {
	general := GeneralCleaningDetails{
		HomeType: generalCleaning.HomeType,
		SQM:      generalCleaning.SQM,
	}
	return json.Marshal(general)
}

func UnmarshalGeneralDetails(detailsOut []byte) (*GeneralCleaningDetails, error) {
	var generalDetails GeneralCleaningDetails
	if err := json.Unmarshal(detailsOut, &generalDetails); err != nil {
		return nil, err
	}
	return &generalDetails, nil
}

func (c CouchCleaningDetails) ToProto() *booking.CouchCleaningDetails {
	specs := make([]*booking.CouchCleaningSpecifications, len(c.CleaningSpecs))
	for i, s := range c.CleaningSpecs {
		specs[i] = &booking.CouchCleaningSpecifications{
			CouchType: booking.CouchType(booking.CouchType_value[s.CouchType]),
			WidthCm:   s.WidthCM,
			DepthCm:   s.DepthCM,
			HeightCm:  s.HeightCM,
			Quantity:  s.Quantity,
		}
	}
	return &booking.CouchCleaningDetails{CleaningSpecs: specs, BedPillows: c.BedPillows}
}

func CouchCleaningDetailsFromProto(pb *booking.CouchCleaningDetails) *CouchCleaningDetails {
	if pb == nil {
		return &CouchCleaningDetails{}
	}
	specs := make([]CouchCleaningSpecifications, len(pb.CleaningSpecs))
	for i, s := range pb.CleaningSpecs {
		specs[i] = CouchCleaningSpecifications{
			CouchType: s.CouchType.String(),
			WidthCM:   s.WidthCm,
			DepthCM:   s.DepthCm,
			HeightCM:  s.HeightCm,
			Quantity:  s.Quantity,
		}
	}
	return &CouchCleaningDetails{CleaningSpecs: specs}
}

func (c *CouchCleaningDetails) MarshalCouchDetails() ([]byte, error) {
	return json.Marshal(c)
}

func UnmarshalCouchDetails(data []byte) (*CouchCleaningDetails, error) {
	var details CouchCleaningDetails
	if err := json.Unmarshal(data, &details); err != nil {
		return nil, err
	}
	return &details, nil
}

// --- Mattress Cleaning ---

func (m MattressCleaningDetails) ToProto() *booking.MattressCleaningDetails {
	specs := make([]*booking.MattressCleaningSpecifications, len(m.CleaningSpecs))
	for i, s := range m.CleaningSpecs {
		specs[i] = &booking.MattressCleaningSpecifications{
			BedType:  booking.BedType(booking.BedType_value[s.BedType]),
			WidthCm:  s.WidthCM,
			DepthCm:  s.DepthCM,
			HeightCm: s.HeightCM,
			Quantity: s.Quantity,
		}
	}
	return &booking.MattressCleaningDetails{CleaningSpecs: specs}
}

func MattressCleaningDetailsFromProto(pb *booking.MattressCleaningDetails) *MattressCleaningDetails {
	if pb == nil {
		return &MattressCleaningDetails{}
	}
	specs := make([]MattressCleaningSpecifications, len(pb.CleaningSpecs))
	for i, s := range pb.CleaningSpecs {
		specs[i] = MattressCleaningSpecifications{
			BedType:  s.BedType.String(),
			WidthCM:  s.WidthCm,
			DepthCM:  s.DepthCm,
			HeightCM: s.HeightCm,
			Quantity: s.Quantity,
		}
	}
	return &MattressCleaningDetails{CleaningSpecs: specs}
}

func (m *MattressCleaningDetails) MarshalMattressDetails() ([]byte, error) {
	return json.Marshal(m)
}

func UnmarshalMattressDetails(data []byte) (*MattressCleaningDetails, error) {
	var details MattressCleaningDetails
	if err := json.Unmarshal(data, &details); err != nil {
		return nil, err
	}
	return &details, nil
}

// --- Car Cleaning ---

func (c CarCleaningDetails) ToProto() *booking.CarCleaningDetails {
	specs := make([]*booking.CarCleaningSpecifications, len(c.CleaningSpecs))
	for i, s := range c.CleaningSpecs {
		specs[i] = &booking.CarCleaningSpecifications{
			CarType:  booking.CarType(booking.CarType_value[s.CarType]),
			Quantity: s.Quantity,
		}
	}
	return &booking.CarCleaningDetails{
		CleaningSpecs: specs,
		ChildSeats:    c.ChildSeats,
	}
}

func CarCleaningDetailsFromProto(pb *booking.CarCleaningDetails) *CarCleaningDetails {
	if pb == nil {
		return &CarCleaningDetails{}
	}
	specs := make([]CarCleaningSpecifications, len(pb.CleaningSpecs))
	for i, s := range pb.CleaningSpecs {
		specs[i] = CarCleaningSpecifications{
			CarType:  s.CarType.String(),
			Quantity: s.Quantity,
		}
	}
	return &CarCleaningDetails{
		CleaningSpecs: specs,
		ChildSeats:    pb.ChildSeats,
	}
}

func (c *CarCleaningDetails) MarshalCarDetails() ([]byte, error) {
	return json.Marshal(c)
}

func UnmarshalCarDetails(data []byte) (*CarCleaningDetails, error) {
	var details CarCleaningDetails
	if err := json.Unmarshal(data, &details); err != nil {
		return nil, err
	}
	return &details, nil
}

func (postConstruction PostConstructionDetails) ToProto() *booking.PostConstructionCleaningDetails {
	return &booking.PostConstructionCleaningDetails{
		Sqm: postConstruction.SQM,
	}
}

func PostConstructionCleaningDetailsFromProto(pb *booking.PostConstructionCleaningDetails) *PostConstructionDetails {
	if pb == nil {
		return &PostConstructionDetails{}
	}
	return &PostConstructionDetails{
		SQM: pb.Sqm,
	}
}

func (postConstructionCleaning *PostConstructionDetails) MarshalPostConstructionDetails() ([]byte, error) {
	postConstruction := PostConstructionDetails{
		SQM: postConstructionCleaning.SQM,
	}
	return json.Marshal(postConstruction)
}

func UnmarshalPostConstructionDetails(detailsOut []byte) (*PostConstructionDetails, error) {
	var postDetails PostConstructionDetails
	if err := json.Unmarshal(detailsOut, &postDetails); err != nil {
		return nil, err
	}
	return &postDetails, nil
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

	if sd.General != (&GeneralCleaningDetails{}) {
		pb.Type = &booking.ServiceDetail_General{
			General: sd.General.ToProto(),
		}
	} else if sd.Couch != (&CouchCleaningDetails{}) {
		pb.Type = &booking.ServiceDetail_Couch{
			Couch: sd.Couch.ToProto(),
		}
	} else if sd.Mattress != (&MattressCleaningDetails{}) {
		pb.Type = &booking.ServiceDetail_Mattress{
			Mattress: sd.Mattress.ToProto(),
		}
	} else if sd.Car != (&CarCleaningDetails{}) {
		pb.Type = &booking.ServiceDetail_Car{
			Car: sd.Car.ToProto(),
		}
	} else if sd.Post != (&PostConstructionDetails{}) {
		pb.Type = &booking.ServiceDetail_Post{
			Post: sd.Post.ToProto(),
		}
	}

	return pb
}
func pointerToTime(t time.Time) *time.Time {
	return &t
}

func FromProtoCreateBooking(req *booking.CreateBookingRequest) *CreateBookingEvent {
	return &CreateBookingEvent{
		Base: BaseBookingDetails{
			CustID:            req.Base.CustId,
			CustomerFirstName: req.Base.CustomerFirstName,
			CustomerLastName:  req.Base.CustomerLastName,
			Address: Address{
				AddressHuman: req.Base.Address.AddressHuman,
				AddressLat:   req.Base.Address.AddressLat,
				AddressLng:   req.Base.Address.AddressLng,
			},
			StartSched:    req.Base.StartSched.AsTime(),
			EndSched:      req.Base.EndSched.AsTime(),
			DirtyScale:    req.Base.DirtyScale,
			PaymentStatus: req.Base.PaymentStatus,
			ReviewStatus:  req.Base.ReviewStatus,
			Photos:        req.Base.Photos,
			CreatedAt:     req.Base.CreatedAt.AsTime(),
			UpdatedAt:     pointerToTime(req.Base.UpdatedAt.AsTime()),
			QuoteId:       req.Base.QuoteId,
		},
		MainService: ServicesRequest{
			ServiceType: MainServiceType(req.MainService.ServiceType.String()),
			Details:     convertServiceDetail(req.MainService.Details),
		},
		Addons: ConvertAddons(req.Addons),
	}
}

func ConvertAddons(addons []*booking.AddOnRequest) []AddOnRequest {
	out := make([]AddOnRequest, 0, len(addons))
	for _, a := range addons {
		out = append(out, AddOnRequest{
			ServiceDetail: ServicesRequest{
				ServiceType: MainServiceType(a.ServiceDetail.ServiceType.String()),
				Details:     convertServiceDetail(a.ServiceDetail.Details),
			},
		})
	}
	return out
}

func convertServiceDetail(sd *booking.ServiceDetail) ServiceDetail {
	if sd == nil {
		return ServiceDetail{}
	}

	switch t := sd.Type.(type) {
	case *booking.ServiceDetail_General:
		return ServiceDetail{
			General: &GeneralCleaningDetails{
				HomeType: t.General.HomeType.String(),
				SQM:      t.General.Sqm,
			},
		}

	case *booking.ServiceDetail_Couch:
		specs := make([]CouchCleaningSpecifications, 0, len(t.Couch.CleaningSpecs))
		for _, s := range t.Couch.CleaningSpecs {
			specs = append(specs, CouchCleaningSpecifications{
				CouchType: s.CouchType.String(),
				WidthCM:   s.WidthCm,
				DepthCM:   s.DepthCm,
				HeightCM:  s.HeightCm,
				Quantity:  s.Quantity,
			})
		}
		return ServiceDetail{
			Couch: &CouchCleaningDetails{CleaningSpecs: specs},
		}

	case *booking.ServiceDetail_Mattress:
		specs := make([]MattressCleaningSpecifications, 0, len(t.Mattress.CleaningSpecs))
		for _, s := range t.Mattress.CleaningSpecs {
			specs = append(specs, MattressCleaningSpecifications{
				BedType:  s.BedType.String(),
				WidthCM:  s.WidthCm,
				DepthCM:  s.DepthCm,
				HeightCM: s.HeightCm,
				Quantity: s.Quantity,
			})
		}
		return ServiceDetail{
			Mattress: &MattressCleaningDetails{CleaningSpecs: specs},
		}

	case *booking.ServiceDetail_Car:
		specs := make([]CarCleaningSpecifications, 0, len(t.Car.CleaningSpecs))
		for _, s := range t.Car.CleaningSpecs {
			specs = append(specs, CarCleaningSpecifications{
				CarType:  s.CarType.String(),
				Quantity: s.Quantity,
			})
		}
		return ServiceDetail{
			Car: &CarCleaningDetails{
				CleaningSpecs: specs,
				ChildSeats:    t.Car.ChildSeats,
			},
		}

	case *booking.ServiceDetail_Post:
		return ServiceDetail{
			Post: &PostConstructionDetails{
				SQM: t.Post.Sqm,
			},
		}

	default:
		return ServiceDetail{}
	}
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
func (cleanerAssigned CleanerAssigned) ToProto() *booking.CleanerAssigned {
	return &booking.CleanerAssigned{
		Id:               cleanerAssigned.ID,
		CleanerFirstName: cleanerAssigned.CleanerFirstName,
		CleanerLastName:  cleanerAssigned.CleanerLastName,
		PfpUrl:           cleanerAssigned.PFPUrl,
	}
}

func CleanerAssignedToProto(cleanerAssignees []CleanerAssigned) []*booking.CleanerAssigned {
	protos := make([]*booking.CleanerAssigned, 0, len(cleanerAssignees))
	for _, cleanerAssigned := range cleanerAssignees {
		protos = append(protos, cleanerAssigned.ToProto())
	}
	return protos
}
func (cleanerEquipment CleaningEquipment) ToProto() *booking.CleaningEquipment {
	return &booking.CleaningEquipment{
		Id:       cleanerEquipment.ID,
		Name:     cleanerEquipment.Name,
		Type:     cleanerEquipment.Type,
		PhotoUrl: cleanerEquipment.PhotoURL,
	}
}

func CleaningEquipmentsToProto(equipments []CleaningEquipment) []*booking.CleaningEquipment {
	protos := make([]*booking.CleaningEquipment, 0, len(equipments))
	for _, equipment := range equipments {
		protos = append(protos, equipment.ToProto())
	}
	return protos
}
func (cleaningResources CleaningResources) ToProto() *booking.CleaningResources {
	return &booking.CleaningResources{
		Id:       cleaningResources.ID,
		Name:     cleaningResources.Name,
		Type:     cleaningResources.Type,
		PhotoUrl: cleaningResources.PhotoURL,
	}
}

func CleaningResourceToProto(resources []CleaningResources) []*booking.CleaningResources {
	protos := make([]*booking.CleaningResources, 0, len(resources))
	for _, resource := range resources {
		protos = append(protos, resource.ToProto())
	}
	return protos
}
func (address Address) ToProto() *booking.Address {
	return &booking.Address{
		AddressHuman: address.AddressHuman,
		AddressLat:   address.AddressLat,
		AddressLng:   address.AddressLng,
	}
}

func AddressFromProto(pb *booking.Address) Address {
	if pb == nil {
		return Address{}
	}
	return Address{
		AddressHuman: pb.AddressHuman,
		AddressLat:   pb.AddressLat,
		AddressLng:   pb.AddressLng,
	}
}
func (baseBooking BaseBookingDetails) ToProto() *booking.BaseBookingDetails {
	var updatedAt *timestamppb.Timestamp
	if baseBooking.UpdatedAt != nil {
		updatedAt = timestamppb.New(*baseBooking.UpdatedAt)
	}

	return &booking.BaseBookingDetails{
		Id:                baseBooking.ID,
		CustId:            baseBooking.CustID,
		CustomerFirstName: baseBooking.CustomerFirstName,
		CustomerLastName:  baseBooking.CustomerLastName,
		Address:           baseBooking.Address.ToProto(),
		StartSched:        timestamppb.New(baseBooking.StartSched),
		EndSched:          timestamppb.New(baseBooking.EndSched),
		DirtyScale:        baseBooking.DirtyScale,
		PaymentStatus:     baseBooking.PaymentStatus,
		ReviewStatus:      baseBooking.ReviewStatus,
		Photos:            baseBooking.Photos,
		CreatedAt:         timestamppb.New(baseBooking.CreatedAt),
		UpdatedAt:         updatedAt,
		QuoteId:           baseBooking.QuoteId,
	}
}

func BaseBookingDetailsFromProto(pb *booking.BaseBookingDetails) BaseBookingDetails {
	if pb == nil {
		return BaseBookingDetails{}
	}

	var updatedAt *time.Time
	if pb.UpdatedAt != nil {
		t := pb.UpdatedAt.AsTime()
		updatedAt = &t
	}

	return BaseBookingDetails{
		ID:                pb.Id,
		CustID:            pb.CustId,
		CustomerFirstName: pb.CustomerFirstName,
		CustomerLastName:  pb.CustomerLastName,
		Address:           AddressFromProto(pb.Address),
		StartSched:        pb.StartSched.AsTime(),
		EndSched:          pb.EndSched.AsTime(),
		DirtyScale:        pb.DirtyScale,
		PaymentStatus:     pb.PaymentStatus,
		ReviewStatus:      pb.ReviewStatus,
		Photos:            pb.Photos,
		CreatedAt:         pb.CreatedAt.AsTime(),
		UpdatedAt:         updatedAt,
		QuoteId:           pb.QuoteId,
	}
}
func (addOns AddOns) ToProto() *booking.AddOn {
	return &booking.AddOn{
		Id:            addOns.ID,
		ServiceDetail: addOns.ServiceDetail.ToProto(),
		Price:         addOns.Price,
	}
}

func AddOnsToProto(addOns []AddOns) []*booking.AddOn {
	protos := make([]*booking.AddOn, 0, len(addOns))
	for _, addOn := range addOns {
		protos = append(protos, addOn.ToProto())
	}
	return protos
}

func (book Booking) ToProto() *booking.Booking {
	return &booking.Booking{
		Id:          book.ID,
		Base:        book.Base.ToProto(),
		MainService: book.MainService.ToProto(),
		Addons:      AddOnsToProto(book.Addons),
		Equipment:   CleaningEquipmentsToProto(book.Equipments),
		Resources:   CleaningResourceToProto(book.Resources),
		Cleaners:    CleanerAssignedToProto(book.Cleaners),
		TotalPrice:  book.TotalPrice,
	}
}
