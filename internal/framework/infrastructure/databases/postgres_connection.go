package databases

import (
	"fmt"
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
	dbClient, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  cnn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}))
	fmt.Println("...cnn: ", cnn)
	fmt.Printf("User: %s, DbName: %s\n", cfg.User, cfg.DbName)

	if err != nil {
		return nil, err
	}

	sqlDB, _ := dbClient.DB()
	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}

	return dbClient, nil
}
