package service

import "handworks/common/grpc/booking"

func CalculateGeneralCleaning(details *booking.GeneralCleaningDetails) float32 {
	sqm := details.Sqm
	homeType := details.HomeType
	if homeType == booking.HomeType_CONDO_ROOM || sqm > 0 && sqm < 31 {
		return 2000.00
	} else if homeType == booking.HomeType_HOUSE || sqm > 31 && sqm < 50 {
		return 2500.00
	} else {
		return 5000.00
	}
}

func CalculateCarCleaning(details *booking.CarCleaningDetails) float32 {
	if details == nil || len(details.CleaningSpecs) == 0 {
		return 0.0
	}

	var total float32 = 0.0

	for _, spec := range details.CleaningSpecs {
		var pricePerCar float32

		switch spec.CarType {
		case booking.CarType_SEDAN:
			pricePerCar = 3250.00
		case booking.CarType_MPV:
			pricePerCar = 4000.00
		case booking.CarType_SUV:
			pricePerCar = 4500.00
		case booking.CarType_VAN:
			pricePerCar = 5000.00
		default:
			pricePerCar = 0.00
		}

		total += pricePerCar * float32(spec.Quantity)
	}
	return total
}

func (b *BookingService) CalculatePriceByServiceType(service *booking.ServicesRequest) float32 {
	switch service.ServiceType {
	case booking.MainServiceType_GENERAL_CLEANING:
		general := service.GetDetails().GetGeneral()
		calculatedPrice := CalculateGeneralCleaning(general)
		return calculatedPrice
	case booking.MainServiceType_COUCH:
	case booking.MainServiceType_MATTRESS:
	case booking.MainServiceType_CAR:
	case booking.MainServiceType_POST:
	default:
		b.L.Error("UNSPECIFIED SERVICE TYPE")
		return 0.0
	}
	return 0.0
}
