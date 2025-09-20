package types

import "handworks/common/grpc/booking"

type CleaningEquipment struct {
	ID       string
	Name     string
	Type     string
	PhotoURL string
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
