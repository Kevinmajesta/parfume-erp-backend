package builder

import (
	"github.com/Kevinmajesta/parfume-erp-backend/internal/entity"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/http/handler"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/http/router"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/repository"
	"github.com/Kevinmajesta/parfume-erp-backend/internal/service"
	"github.com/Kevinmajesta/parfume-erp-backend/pkg/cache"
	"github.com/Kevinmajesta/parfume-erp-backend/pkg/email"
	"github.com/Kevinmajesta/parfume-erp-backend/pkg/encrypt"
	"github.com/Kevinmajesta/parfume-erp-backend/pkg/route"

	// "github.com/labstack/echo/"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func BuildPublicRoutes(db *gorm.DB, redisDB *redis.Client, entityCfg *entity.Config, encryptTool encrypt.EncryptTool) []*route.Route {
	cacheable := cache.NewCacheable(redisDB)
	emailService := email.NewEmailSender(entityCfg)
	userRepository := repository.NewUserRepository(db, cacheable)
	userService := service.NewUserService(userRepository, encryptTool, emailService)
	userHandler := handler.NewUserHandler(userService)

	adminRepository := repository.NewAdminRepository(db, cacheable)
	adminService := service.NewAdminService(adminRepository, encryptTool, emailService)
	adminHandler := handler.NewAdminHandler(adminService)

	return router.PublicRoutes(userHandler, adminHandler)
}

func BuildPrivateRoutes(db *gorm.DB, redisDB *redis.Client, encryptTool encrypt.EncryptTool, entityCfg *entity.Config) []*route.Route {
	cacheable := cache.NewCacheable(redisDB)
	emailService := email.NewEmailSender(entityCfg)
	userRepository := repository.NewUserRepository(db, cacheable)
	userService := service.NewUserService(userRepository, encryptTool, nil)
	userHandler := handler.NewUserHandler(userService)

	suggestionRepository := repository.NewSuggestionRepository(db, cacheable)
	suggestionService := service.NewSuggestionService(suggestionRepository, userRepository)
	suggestionHandler := handler.NewSuggestionHandler(suggestionService, userService)

	adminRepository := repository.NewAdminRepository(db, cacheable)
	adminService := service.NewAdminService(adminRepository, encryptTool, nil)
	adminHandler := handler.NewAdminHandler(adminService)

	schedulesRepository := repository.NewSchedulesRepository(db, cacheable)
	schedulesService := service.NewSchedulesService(schedulesRepository)
	schedulesHandler := handler.NewSchedulesHandler(schedulesService)

	productRepository := repository.NewProductRepository(db, cacheable)
	productService := service.NewProductService(productRepository)
	productHandler := handler.NewProductHandler(productService)

	materialRepository := repository.NewMaterialRepository(db, cacheable)
	materialService := service.NewMaterialService(materialRepository)
	materialHandler := handler.NewMaterialHandler(materialService)

	bomRepository := repository.NewBOMRepository(db, cacheable)
	bomMaterialRepo := repository.NewBOMMaterialRepository(db)
	bomService := service.NewBOMService(bomRepository, bomMaterialRepo)
	bomHandler := handler.NewBOMHandler(bomService)

	moRepository := repository.NewMoRepository(db, cacheable)
	moService := service.NewMoService(moRepository)
	moHandler := handler.NewMoHandler(moService)

	vendorRepository := repository.NewVendorRepository(db, cacheable)
	vendorService := service.NewVendorService(vendorRepository)
	vendorHandler := handler.NewVendorHandler(vendorService)

	rfqRepository := repository.NewRfqRepository(db, cacheable)
	rfqProductRepo := repository.NewRfqProductRepository(db)
	rfqService := service.NewRfqService(rfqRepository, rfqProductRepo, emailService)
	rfqHandler := handler.NewRfqHandler(rfqService)

	costumerRepository := repository.NewCostumerRepository(db, cacheable)
	costumerService := service.NewCostumerService(costumerRepository)
	costumerHandler := handler.NewCostumerHandler(costumerService)

	quoRepository := repository.NewQuoRepository(db, cacheable)
	quoProductRepo := repository.NewQuoProductRepository(db)
	quoService := service.NewQuoService(quoRepository, quoProductRepo, emailService)
	quoHandler := handler.NewQuoHandler(quoService)

	return router.PrivateRoutes(userHandler, suggestionHandler, adminHandler, schedulesHandler,
		productHandler, materialHandler, *bomHandler, moHandler, vendorHandler, rfqHandler, costumerHandler, quoHandler)
}
