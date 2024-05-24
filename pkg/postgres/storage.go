package postgres

import (
	"EffectiveMobile/config"
	"EffectiveMobile/pkg/postgres/models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
)

type Storage struct {
	Db *gorm.DB
}

var (
	once       sync.Once
	dbInstance *Storage
)

func NewStorage(conf *config.Config) *Storage {
	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			conf.PSQL.Host,
			conf.PSQL.User,
			conf.PSQL.Password,
			conf.PSQL.DBName,
			conf.PSQL.Port,
			conf.PSQL.SSLMode,
		)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		err = db.AutoMigrate(&models.Car{}, &models.People{})
		if err != nil {
			return
		}

		dbInstance = &Storage{Db: db}
	})
	return dbInstance
}
