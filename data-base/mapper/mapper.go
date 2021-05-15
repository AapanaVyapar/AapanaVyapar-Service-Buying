package mapper

import (
	"aapanavyapar-service-buying/data-base/structs"
	"time"
)

func MapLocationToSector(location *structs.Location) int32 {

	return 10
}

func CalculateDeliveryTime(distance int) time.Time {

	return time.Now().UTC()
}

func CalculateDeliveryCost(distance int, address *structs.Address) float32 {

	return float32(distance) * 5
}
