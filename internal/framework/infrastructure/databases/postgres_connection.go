package databases

import (
	"fmt"
	"log"
	"time"

	"github.com/ali-mahdavi-dev/bunny-go/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(cfg config.PostgresConfig) (*gorm.DB, error) {
	var err error
	cnn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Tehran",
		cfg.Host, cfg.Port, cfg.User, cfg.Password,
		cfg.DbName, cfg.SSLMode)

	dbClient, err := gorm.Open(postgres.Open(cnn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDb, _ := dbClient.DB()
	err = sqlDb.Ping()
	if err != nil {
		return nil, err
	}

	sqlDb.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDb.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDb.SetConnMaxLifetime(cfg.ConnMaxLifetime * time.Minute)

	log.Println("Db connection established")

	return dbClient, nil
}
