package builder

import (
	"github.com/Kevinmajesta/webPemancingan/internal/entity"
	"github.com/Kevinmajesta/webPemancingan/internal/http/handler"
	"github.com/Kevinmajesta/webPemancingan/internal/http/router"
	"github.com/Kevinmajesta/webPemancingan/internal/repository"
	"github.com/Kevinmajesta/webPemancingan/internal/service"
	"github.com/Kevinmajesta/webPemancingan/pkg/cache"
	"github.com/Kevinmajesta/webPemancingan/pkg/email"
	"github.com/Kevinmajesta/webPemancingan/pkg/encrypt"
	"github.com/Kevinmajesta/webPemancingan/pkg/route"
	"github.com/Kevinmajesta/webPemancingan/pkg/token"

	// "github.com/labstack/echo/"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func BuildPublicRoutes(db *gorm.DB, redisDB *redis.Client, entityCfg *entity.Config, tokenUseCase token.TokenUseCase, encryptTool encrypt.EncryptTool) []*route.Route {
	cacheable := cache.NewCacheable(redisDB)
	emailService := email.NewEmailSender(entityCfg)
	userRepository := repository.NewUserRepository(db, cacheable)
	userService := service.NewUserService(userRepository, tokenUseCase, encryptTool, emailService)
	userHandler := handler.NewUserHandler(userService)

	adminRepository := repository.NewAdminRepository(db, cacheable)
	adminService := service.NewAdminService(adminRepository, tokenUseCase, encryptTool, emailService)
	adminHandler := handler.NewAdminHandler(adminService)

	return router.PublicRoutes(userHandler, adminHandler)
}

func BuildPrivateRoutes(db *gorm.DB, redisDB *redis.Client, encryptTool encrypt.EncryptTool, tokenUseCase token.TokenUseCase) []*route.Route {
	cacheable := cache.NewCacheable(redisDB)
	userRepository := repository.NewUserRepository(db, cacheable)
	userService := service.NewUserService(userRepository, nil, encryptTool, nil)
	userHandler := handler.NewUserHandler(userService)

	suggestionRepository := repository.NewSuggestionRepository(db, cacheable)
	suggestionService := service.NewSuggestionService(suggestionRepository, tokenUseCase, userRepository)
	suggestionHandler := handler.NewSuggestionHandler(suggestionService, userService)

	adminRepository := repository.NewAdminRepository(db, cacheable)
	adminService := service.NewAdminService(adminRepository, tokenUseCase, encryptTool, nil)
	adminHandler := handler.NewAdminHandler(adminService)

	schedulesRepository := repository.NewSchedulesRepository(db, cacheable)
	schedulesService := service.NewSchedulesService(schedulesRepository)
	schedulesHandler := handler.NewSchedulesHandler(schedulesService)

	return router.PrivateRoutes(userHandler, suggestionHandler, adminHandler, schedulesHandler)
}
