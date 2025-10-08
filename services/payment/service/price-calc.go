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
	return total
}

func CalculateCouchCleaning(details *booking.CouchCleaningDetails) float32 {
	if details == nil || len(details.CleaningSpecs) == 0 {
		return 0.0
	}
	var total float32 = 0.0

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
			pricePerCouch = 500
		case booking.CouchType_LAZBOY:
			pricePerCouch = 900
		case booking.CouchType_CHAIR:
			pricePerCouch = 250
		}
		total += pricePerCouch * float32(spec.Quantity)
	}
	return total
}

func CalculatePostConstructionCleaning(details *booking.PostConstructionCleaningDetails) float32 {
	total := details.Sqm * 50
	return float32(total)
}
func (b *PaymentService) CalculatePriceByServiceType(service *booking.ServicesRequest) float32 {
	var calculatedPrice float32 = 0.00
	switch service.ServiceType {
	case booking.MainServiceType_GENERAL_CLEANING:
		general := service.GetDetails().GetGeneral()
		calculatedPrice = CalculateGeneralCleaning(general)
		return calculatedPrice
	case booking.MainServiceType_COUCH:
	case booking.MainServiceType_MATTRESS:
	case booking.MainServiceType_CAR:
		car := service.GetDetails().GetCar()
		calculatedPrice = CalculateCarCleaning(car)
		return calculatedPrice
	case booking.MainServiceType_POST:
	default:
		b.L.Error("UNSPECIFIED SERVICE TYPE")
		return calculatedPrice
	}
	return calculatedPrice
}
