package component

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"rbd_proxy_dp/config"
	"rbd_proxy_dp/model"
)

type DB struct {
	db *gorm.DB
}

var db *gorm.DB

func NewDB() *DB {
	return &DB{
		db: db,
	}
}

func (d *DB) Start(ctx context.Context) error {
	dbConfig := config.DefaultDB()

	var err error
	switch dbConfig.DBType {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			dbConfig.DBUser, dbConfig.DBPasswd, dbConfig.DBHost, dbConfig.DBPort, dbConfig.DBName)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			logrus.Errorf("failed to connect to mysql: %v", err)
			return fmt.Errorf("mysql connection error: %w", err)
		}
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dbConfig.DBName), &gorm.Config{})
		if err != nil {
			logrus.Errorf("failed to connect to sqlite: %v", err)
			return fmt.Errorf("sqlite connection error: %w", err)
		}
	default:
		logrus.Errorf("unsupported db_type: %s", dbConfig.DBType)
		return fmt.Errorf("unsupported db_type: %s", dbConfig.DBType)
	}

	// 自动迁移模型
	if err := db.AutoMigrate(&model.APIResponse{}); err != nil {
		logrus.Errorf("failed to migrate database: %v", err)
		return fmt.Errorf("database migration error: %w", err)
	}

	logrus.Infof("database connection established and migrated.")
	return nil
}

func (d *DB) CloseHandle() {

}
