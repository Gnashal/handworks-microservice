package service

import "handworks/common/grpc/booking"

// hash maps
var mattressPrices = map[booking.BedType]float32{
	booking.BedType_KING:           2000.00,
	booking.BedType_KING_HEADBAND:  2500.00,
	booking.BedType_QUEEN:          1800.00,
	booking.BedType_QUEEN_HEADBAND: 2300.00,
	booking.BedType_DOUBLE:         1500.00,
	booking.BedType_SINGLE:         1000.00,
}
var carPrices = map[booking.CarType]float32{
	booking.CarType_SEDAN:     3250.00,
	booking.CarType_MPV:       4000.00,
	booking.CarType_SUV:       4000.00,
	booking.CarType_VAN:       5200.00,
	booking.CarType_PICKUP:    3600.00,
	booking.CarType_CAR_SMALL: 1750.00,
}

var couchPrices = map[booking.CouchType]float32{
	booking.CouchType_SEATER_1:             500.00,
	booking.CouchType_SEATER_2:             1000.00,
	booking.CouchType_SEATER_3:             1300.00,
	booking.CouchType_SEATER_3_LTYPE_SMALL: 1500.00,
	booking.CouchType_SEATER_3_LTYPE_LARGE: 1750.00,
	booking.CouchType_SEATER_4_LTYPE_SMALL: 1800.00,
	booking.CouchType_SEATER_4_LTYPE_LARGE: 2000.00,
	booking.CouchType_SEATER_5_LTYPE:       2250.00,
	booking.CouchType_SEATER_6_LTYPE:       2500.00,
	booking.CouchType_OTTOMAN:              500.00,
	booking.CouchType_LAZBOY:               900.00,
	booking.CouchType_CHAIR:                250.00,
}

func CalculateGeneralCleaning(details *booking.GeneralCleaningDetails) float32 {
	if details == nil {
		return 0.0
	}

	sqm := details.Sqm
	homeType := details.HomeType

	switch {
	case homeType == booking.HomeType_CONDO_ROOM || (sqm > 0 && sqm <= 30):
		return 2000.00
	case homeType == booking.HomeType_HOUSE || (sqm > 30 && sqm <= 50):
		return 2500.00
	case sqm > 50 && sqm <= 100:
		return 5000.00
	default:
		return float32(sqm * 50)
	}
}

func CalculateCarCleaning(details *booking.CarCleaningDetails) float32 {
	if details == nil {
		return 0.0
	}

	var total float32
	for _, spec := range details.CleaningSpecs {
		price := carPrices[spec.CarType]
		total += price * float32(spec.Quantity)
	}

	if details.ChildSeats > 0 {
		total += float32(details.ChildSeats) * 250.00
	}

	return total
}

func CalculateCouchCleaning(details *booking.CouchCleaningDetails) float32 {
	if details == nil {
		return 0.0
	}

	var total float32
	for _, spec := range details.CleaningSpecs {
		price := couchPrices[spec.CouchType]
		total += price * float32(spec.Quantity)
	}

	if details.BedPillows > 0 {
		total += float32(details.BedPillows) * 100.00
	}

	return total
}

func CalculateMattressCleaning(details *booking.MattressCleaningDetails) float32 {
	if details == nil {
		return 0.0
	}

	var total float32
	for _, spec := range details.CleaningSpecs {
		price := mattressPrices[spec.BedType]
		total += price * float32(spec.Quantity)
	}
	return total
}

func CalculatePostConstructionCleaning(details *booking.PostConstructionCleaningDetails) float32 {
	if details == nil {
		return 0.0
	}
	return float32(details.Sqm * 50.00)
}

func (b *PaymentService) CalculatePriceByServiceType(service *booking.ServicesRequest) float32 {
	if service == nil {
		b.L.Error("service is nil")
		return 0
	}

	if service.GetDetails() == nil {
		b.L.Error("service.Details is nil for service type: %v", service.ServiceType)
		return 0
	}
	var calculatedPrice float32 = 0.00

	switch service.ServiceType {
	case booking.MainServiceType_GENERAL_CLEANING:
		calculatedPrice = CalculateGeneralCleaning(service.GetDetails().GetGeneral())

	case booking.MainServiceType_COUCH:
		calculatedPrice = CalculateCouchCleaning(service.GetDetails().GetCouch())

	case booking.MainServiceType_MATTRESS:
		calculatedPrice = CalculateMattressCleaning(service.GetDetails().GetMattress())

	case booking.MainServiceType_CAR:
		calculatedPrice = CalculateCarCleaning(service.GetDetails().GetCar())

	case booking.MainServiceType_POST:
		calculatedPrice = CalculatePostConstructionCleaning(service.GetDetails().GetPost())

	default:
		b.L.Error("UNSPECIFIED SERVICE TYPE")
	}

	return calculatedPrice
}
