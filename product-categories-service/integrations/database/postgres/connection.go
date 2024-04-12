package postgres

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/url"
	"product-categories-service/config"
	"product-categories-service/domains"
	"product-categories-service/utils"
)

var Connection *gorm.DB

func SetupDatabaseConnection() {
	cfg := config.Configuration
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		url.QueryEscape(cfg.Database.User),
		url.QueryEscape(cfg.Database.Password),
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DatabaseName,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: utils.GormLogger{},
	})
	if err != nil {
		panic("failed to create a connection to database")
	}
	db.AutoMigrate(
		&domains.Category{},
	)
	Connection = db
}
func CloseDatabaseConnection() {
	dbSQL, err := Connection.DB()
	if err != nil {
		panic("failed to close connection from database")
	}
	err = dbSQL.Close()
	if err != nil {
		panic("failed to close connection from database")
	}
}
