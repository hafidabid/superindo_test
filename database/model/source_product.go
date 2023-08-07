package model

type SourceProduct struct {
	Product
}

func (sourceProduct SourceProduct) TableName() string {
	return "source_product"
}
