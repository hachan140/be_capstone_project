package storage

import (
	"be-capstone-project/src/internal/core/common_configs"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

func NewPostgresClient(config *common_configs.Store) (*Database, error) {
	dsn := fmt.Sprintf("postgres://%v:%v@%v:%v/%s", config.User, config.Password, config.Host, config.Port, config.Db)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Cấu hình pool kết nối
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(int(config.MaxIdleConns))

	// Debug mode
	if config.IsDebug {
		db = db.Debug()
	}

	return &Database{db}, nil
}
