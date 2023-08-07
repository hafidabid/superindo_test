package database

import (
	"errors"
	"fmt"
	cliLog "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"superindo_diksha_test/config"
	"superindo_diksha_test/database/model"
	"time"
)

type AppDatabase struct {
	db1 *gorm.DB
	db2 *gorm.DB
}

func NewAppDatabase(dbConfig1, dbConfig2 config.DbConfig) (*AppDatabase, error) {
	lgr := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		})

	gormSetting := gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		Logger:                 lgr,
	}

	db1Uri := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		dbConfig1.PostgresUser,
		dbConfig1.PostgresPassword,
		dbConfig1.PostgresHost,
		dbConfig1.PostgresPort,
		dbConfig1.PostgresDb)
	db2Uri := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		dbConfig2.PostgresUser,
		dbConfig2.PostgresPassword,
		dbConfig2.PostgresHost,
		dbConfig2.PostgresPort,
		dbConfig2.PostgresDb)

	db1, err := gorm.Open(postgres.Open(db1Uri), &gormSetting)
	if err != nil {
		cliLog.Errorf("Error on open connection db1: %v", err)
		return nil, err
	}

	db2, err := gorm.Open(postgres.Open(db2Uri), &gormSetting)
	if err != nil {
		cliLog.Errorf("Error on open connection db2: %v", err)
		return nil, err
	}

	return &AppDatabase{
		db1: db1,
		db2: db2,
	}, nil
}

func (db AppDatabase) Migrate() error {
	if err := db.db1.AutoMigrate(model.SourceProduct{}); err != nil {
		return err
	}
	if err := db.db2.AutoMigrate(model.DestinationProduct{}); err != nil {
		return err
	}
	return nil
}

func (db AppDatabase) Flush() error {
	if err := db.db1.Where("id is not null").Delete(&model.SourceProduct{}).Error; err != nil {
		return err
	}
	if err := db.db2.Where("id is not null").Delete(&model.DestinationProduct{}).Error; err != nil {
		return err
	}
	return nil
}

func (db AppDatabase) InsertProducts(sourceData []model.SourceProduct, destData []model.DestinationProduct) error {
	txDb1 := db.db1.Begin()
	txDb2 := db.db2.Begin()

	if err := txDb1.Debug().Create(sourceData).Error; err != nil {
		txDb1.Rollback()
		return err
	}

	if err := txDb2.Debug().Create(destData).Error; err != nil {
		txDb2.Rollback()
		txDb1.Rollback()
		return err
	}

	if e1 := txDb1.Commit().Error; e1 != nil {
		return e1
	}

	if e2 := txDb2.Commit().Error; e2 != nil {
		return e2
	}

	return nil
}

func (db AppDatabase) CheckUpdateProduct() error {
	dest, err := db.GetSourceProduct(0, 99999999999999999)
	if err != nil {
		return err
	}

	for _, pd := range dest {
		db.db2.Model(model.DestinationProduct{}).Where("id = ?", pd.ID).Updates(map[string]interface{}{
			"updated_at":    time.Now(),
			"qty":           pd.Qty,
			"selling_price": pd.SellingPrice,
			"promo_price":   pd.PromoPrice,
		})
	}
	return nil
}

func (db AppDatabase) GetSourceProduct(page int, limit int) ([]model.SourceProduct, error) {
	var products []model.SourceProduct

	err := db.db1.Model(model.SourceProduct{}).Offset(page * limit).Limit(limit).Find(&products).Error
	if err != nil {
		return []model.SourceProduct{}, err
	}

	return products, nil
}

func (db AppDatabase) GetDestinationProduct(page int, limit int) ([]model.DestinationProduct, error) {
	var products []model.DestinationProduct

	err := db.db2.Model(model.DestinationProduct{}).Offset(page * limit).Limit(limit).Find(&products).Error
	if err != nil {
		return []model.DestinationProduct{}, err
	}
	return products, nil
}

func (db AppDatabase) CountProduct(dbId uint) (int64, error) {
	if dbId > 1 {
		return 0, errors.New("dbId should 0/1")
	}

	var mdl any
	engine := db.db1
	mdl = model.SourceProduct{}
	if dbId == 1 {
		engine = db.db2
		mdl = model.DestinationProduct{}
	}

	var count int64
	err := engine.Model(mdl).Count(&count).Error
	if err != nil {
		return 0, nil
	}

	return count, nil
}
