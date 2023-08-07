package controller

import (
	log "github.com/sirupsen/logrus"
	"math"
	"superindo_diksha_test/config"
	"superindo_diksha_test/database"
	"superindo_diksha_test/database/model"
	"superindo_diksha_test/utils"
)

type Controller interface {
}

type ProductController struct {
	Controller
	config config.Config
	db     database.AppDatabase
}

func NewController(conf config.Config, db database.AppDatabase) ProductController {
	return ProductController{
		config: conf,
		db:     db,
	}
}

func (c ProductController) GenerateData(numOfData int) error {
	sourceData, err := utils.GenerateMock(numOfData)
	if err != nil {
		return err
	}

	destData := utils.SanitizeMock(sourceData)

	sourceInsert := []model.SourceProduct{}
	destInsert := []model.DestinationProduct{}

	for _, d := range sourceData {
		sourceInsert = append(sourceInsert, model.SourceProduct{
			Product: d,
		})
	}

	for _, d := range destData {
		destInsert = append(destInsert, model.DestinationProduct{
			Product: d,
		})
	}

	err = c.db.InsertProducts(sourceInsert, destInsert)

	if err != nil {
		return err
	}

	return nil
}

func (c ProductController) GetSource(page int, limit int) ([]model.SourceProduct, *utils.Metadata, error) {
	// get data
	data, err := c.db.GetSourceProduct(page, limit)
	if err != nil {
		log.Errorf("Error in get source product from db with error: %v", err)
		return []model.SourceProduct{}, nil, nil
	}
	// fetch metadata
	counter, err := c.db.CountProduct(0)
	if err != nil {
		log.Errorf("Error in count from db with error: %v", err)
		return []model.SourceProduct{}, nil, nil
	}

	metadata := utils.Metadata{
		Page:      page,
		Limit:     limit,
		TotalPage: math.Ceil(float64(counter) / float64(limit)),
		TotalData: counter,
	}

	return data, &metadata, nil
}

func (c ProductController) GetDestination(page int, limit int) ([]model.DestinationProduct, *utils.Metadata, error) {
	// get data
	data, err := c.db.GetDestinationProduct(page, limit)
	if err != nil {
		log.Errorf("Error in get destination product from db with error: %v", err)
		return []model.DestinationProduct{}, nil, nil
	}

	//fetch metadata
	counter, err := c.db.CountProduct(1)
	if err != nil {
		log.Errorf("Error in count from db with error: %v", err)
		return []model.DestinationProduct{}, nil, nil
	}

	metadata := utils.Metadata{
		Page:      page,
		Limit:     limit,
		TotalPage: math.Ceil(float64(counter) / float64(limit)),
		TotalData: counter,
	}

	return data, &metadata, nil
}

func (c ProductController) UpdateData() error {
	err := c.db.CheckUpdateProduct()
	if err != nil {
		log.Errorf("Error in update data from db with error: %v", err)
		return err
	}
	return nil
}
