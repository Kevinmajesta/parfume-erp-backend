package main

import (
	"github.com/Kevinmajesta/parfume-erp-backend/configs"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/builder"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/pkg/cache"
	"github.com/Kevinmajesta/parfume-erp-backend/pkg/encrypt"
	"github.com/Kevinmajesta/parfume-erp-backend/pkg/postgres"
	"github.com/Kevinmajesta/parfume-erp-backend/pkg/server"
)

func main() {
	// Load configurations from .env file
	cfg, err := configs.NewConfig(".env")
	checkError(err)

	// Initialize PostgreSQL database connection
	db, err := postgres.InitPostgres(&cfg.Postgres)
	checkError(err)

	// Initialize Redis cache connection
	redisDB := cache.InitCache(&cfg.Redis)

	// Initialize encryption tool
	encryptTool := encrypt.NewEncryptTool(cfg.Encrypt.SecretKey, cfg.Encrypt.IV)

	// Convert configs.Config to *entity.Config
	entityCfg := convertToEntityConfig(cfg)

	// Build public and private routes
	publicRoutes := builder.BuildPublicRoutes(db, redisDB, entityCfg , encryptTool)
	privateRoutes := builder.BuildPrivateRoutes(db, redisDB, encryptTool)

	// Initialize and run the server
	srv := server.NewServer("app", publicRoutes, privateRoutes)
	srv.Run()
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// Example function to convert configs.Config to *entity.Config
func convertToEntityConfig(cfg *configs.Config) *entity.Config {
	return &entity.Config{
		SMTP: entity.SMTPConfig{
			Host:     cfg.SMTP.Host,
			Port:     cfg.SMTP.Port,
			Password: cfg.SMTP.Password,
		},
		// Add other fields as needed
	}
}
