package utils

import (
	"github.com/bxcodec/faker/v3"
	"superindo_diksha_test/database/model"
)

func GenerateMock(numOfData int) ([]model.Product, error) {
	var data = make(chan model.Product)
	for i := 0; i < numOfData; i++ {
		i := i
		go func() {
			qty, _ := faker.RandomInt(1, 100)
			sellPrice, _ := faker.RandomInt(1000, 10000)
			promoPrice, _ := faker.RandomInt(500, 9000)

			data <- model.Product{
				ID:           i + 1,
				ProductName:  faker.Username(),
				Qty:          uint(qty[0]),
				SellingPrice: float64(sellPrice[0]),
				PromoPrice:   float64(promoPrice[0]),
			}
		}()
	}

	retData := make([]model.Product, 0, numOfData)

	for i := 0; i < numOfData; i++ {
		retData = append(retData, <-data)
	}
	return retData, nil
}

func SanitizeMock(data []model.Product) []model.Product {
	newData := []model.Product{}

	for _, v := range data {
		v.Qty = 0
		v.SellingPrice = 0
		v.PromoPrice = 0
		newData = append(newData, v)
	}

	return newData
}
