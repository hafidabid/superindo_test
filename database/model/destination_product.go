package model

type DestinationProduct struct {
	Product
}

func (destinationProduct DestinationProduct) TableName() string {
	return "destination_product"
}
