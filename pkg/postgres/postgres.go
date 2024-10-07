package postgres

import (
	"fmt"

	"github.com/Kevinmajesta/webPemancingan/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitPostgres(config *configs.PostgresConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		`host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta`,
		config.Host,
		config.User,
		config.Password,
		config.Database,
		config.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return db, err
	}
	return db, nil
}
