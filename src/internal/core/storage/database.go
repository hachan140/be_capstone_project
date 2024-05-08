package storage

import (
	"be-capstone-project/src/internal/core/common_configs"
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

// Database is embed type for postgres client, only expose interface
type Database struct {
	*bun.DB
}

func NewPostgresClient(config *common_configs.Store) (*Database, error) {
	dsn := fmt.Sprintf("postgres://%v:%v@%v:%v/%s?sslmode=disable", config.User, config.Password, config.Host, config.Port, config.Db)

	sqlDB := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(dsn),
	))

	sqlDB.SetMaxIdleConns(int(config.MaxIdleConns))
	//sqlDB.SetConnMaxLifetime(time.Duration(config.Store.MaxIdleConns))

	if config.IsDebug {
		bundebug.NewQueryHook(bundebug.WithVerbose(true))
	}
	db := bun.NewDB(sqlDB, pgdialect.New())
	err := db.Ping()
	if err != nil {
		return nil, err
	}
	return &Database{db}, nil
}
