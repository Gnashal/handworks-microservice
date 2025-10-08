package service

import "handworks/common/grpc/booking"

func CalculateGeneralCleaning(details *booking.GeneralCleaningDetails) float32 {
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
		var pricePerCar float32

		switch spec.CarType {
		case booking.CarType_SEDAN:
			pricePerCar = 3250.00
		case booking.CarType_MPV:
			pricePerCar = 4000.00
		case booking.CarType_SUV:
			pricePerCar = 4000.00
		case booking.CarType_VAN:
			pricePerCar = 5200.00
		case booking.CarType_PICKUP:
			pricePerCar = 3600.00
		case booking.CarType_CAR_SMALL:
			pricePerCar = 1750.00
		default:
			pricePerCar = 0.00
		}

		total += pricePerCar * float32(spec.Quantity)
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
		var pricePerCouch float32

		switch spec.CouchType {
		case booking.CouchType_SEATER_1:
			pricePerCouch = 500.00
		case booking.CouchType_SEATER_2:
			pricePerCouch = 1000.00
		case booking.CouchType_SEATER_3:
			pricePerCouch = 1300.00
		case booking.CouchType_SEATER_3_LTYPE_SMALL:
			pricePerCouch = 1500.00
		case booking.CouchType_SEATER_3_LTYPE_LARGE:
			pricePerCouch = 1750.00
		case booking.CouchType_SEATER_4_LTYPE_SMALL:
			pricePerCouch = 1800.00
		case booking.CouchType_SEATER_4_LTYPE_LARGE:
			pricePerCouch = 2000.00
		case booking.CouchType_SEATER_5_LTYPE:
			pricePerCouch = 2250.00
		case booking.CouchType_SEATER_6_LTYPE:
			pricePerCouch = 2500.00
		case booking.CouchType_OTTOMAN:
			pricePerCouch = 500.00
		case booking.CouchType_LAZBOY:
			pricePerCouch = 900.00
		case booking.CouchType_CHAIR:
			pricePerCouch = 250.00
		default:
			pricePerCouch = 0.00
		}

		total += pricePerCouch * float32(spec.Quantity)
	}
	if details.BedPillows > 0 {
		total += float32(details.BedPillows) * 100.00
	}

	return total
}

func CalculateMattressCleaning(details *booking.MattressCleaningDetails) float32 {
	if details == nil || len(details.CleaningSpecs) == 0 {
		return 0.0
	}

	var total float32
	for _, spec := range details.CleaningSpecs {
		var price float32

		switch spec.BedType {
		case booking.BedType_KING:
			price = 2000.00
		case booking.BedType_KING_HEADBAND:
			price = 2500.00
		case booking.BedType_QUEEN:
			price = 1800.00
		case booking.BedType_QUEEN_HEADBAND:
			price = 2300.00
		case booking.BedType_DOUBLE:
			price = 1500.00
		case booking.BedType_SINGLE:
			price = 1000.00
		default:
			price = 0.00
		}

		total += price * float32(spec.Quantity)
	}
	return total
}

func CalculatePostConstructionCleaning(details *booking.PostConstructionCleaningDetails) float32 {
	if details == nil {
		return 0.0
	}
	sqm := details.Sqm
	return float32(sqm * 50.00)
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
