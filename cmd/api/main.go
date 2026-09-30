package main

import (
	"log"

	_ "github.com/liyansasongko/bumdes-be/docs"
	"github.com/liyansasongko/bumdes-be/internal/config"
	"github.com/liyansasongko/bumdes-be/internal/database"
	"github.com/liyansasongko/bumdes-be/internal/handler"
	"github.com/liyansasongko/bumdes-be/internal/repository"
	"github.com/liyansasongko/bumdes-be/internal/router"
	"github.com/liyansasongko/bumdes-be/internal/service"
)

// @title BUMDes Backend API
// @version 1.0
// @description REST API untuk aplikasi BUMDes: autentikasi, profil user, manajemen role & permission, dan manajemen akun.
// @basePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Masukkan token dengan format: Bearer {token}
func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("gagal terhubung ke database: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)
	userLogRepo := repository.NewUserLogRepository(db)
	parkingPriceRepo := repository.NewParkingPriceRepository(db)
	paymentMethodRepo := repository.NewPaymentMethodRepository(db)
	revenueCategoryRepo := repository.NewRevenueCategoryRepository(db)
	ticketPriceRepo := repository.NewTicketPriceRepository(db)
	tenantTypeRepo := repository.NewTenantTypeRepository(db)
	attractionPriceRepo := repository.NewAttractionPriceRepository(db)
	parkingRepo := repository.NewParkingRepository(db)
	kiosLocationRepo := repository.NewKiosLocationRepository(db)
	ticketingRepo := repository.NewTicketingRepository(db)
	attractionRepo := repository.NewAttractionRepository(db)
	tenantRepo := repository.NewTenantRepository(db)
	billingHistoryRepo := repository.NewBillingHistoryRepository(db)
	reportRepo := repository.NewReportRepository(db)

	authService := service.NewAuthService(userRepo, cfg)
	userService := service.NewUserService(userRepo)
	roleService := service.NewRoleService(roleRepo, permissionRepo)
	permissionService := service.NewPermissionService(permissionRepo)
	userLogService := service.NewUserLogService(userLogRepo)
	parkingPriceService := service.NewParkingPriceService(parkingPriceRepo)
	paymentMethodService := service.NewPaymentMethodService(paymentMethodRepo)
	revenueCategoryService := service.NewRevenueCategoryService(revenueCategoryRepo)
	ticketPriceService := service.NewTicketPriceService(ticketPriceRepo)
	tenantTypeService := service.NewTenantTypeService(tenantTypeRepo)
	attractionPriceService := service.NewAttractionPriceService(attractionPriceRepo)
	parkingService := service.NewParkingService(parkingRepo, parkingPriceRepo, paymentMethodRepo)
	kiosLocationService := service.NewKiosLocationService(kiosLocationRepo)
	ticketingService := service.NewTicketingService(ticketingRepo, ticketPriceRepo, paymentMethodRepo)
	attractionService := service.NewAttractionService(attractionRepo, attractionPriceRepo, paymentMethodRepo)
	tenantService := service.NewTenantService(tenantRepo, tenantTypeRepo, kiosLocationRepo)
	kiosAvailableService := service.NewKiosAvailableService(kiosLocationRepo, tenantRepo)
	billingService := service.NewBillingService(billingHistoryRepo, tenantRepo, cfg.BillingDueReminderDays)
	reportService := service.NewReportService(reportRepo)

	h := router.Handlers{
		AuthHandler:            handler.NewAuthHandler(authService),
		UserHandler:            handler.NewUserHandler(userService, userLogService),
		RoleHandler:            handler.NewRoleHandler(roleService, userLogService),
		PermissionHandler:      handler.NewPermissionHandler(permissionService),
		UserLogHandler:         handler.NewUserLogHandler(userLogService),
		ParkingPriceHandler:    handler.NewParkingPriceHandler(parkingPriceService, userLogService),
		PaymentMethodHandler:   handler.NewPaymentMethodHandler(paymentMethodService, userLogService),
		RevenueCategoryHandler: handler.NewRevenueCategoryHandler(revenueCategoryService, userLogService),
		TicketPriceHandler:     handler.NewTicketPriceHandler(ticketPriceService, userLogService),
		TenantTypeHandler:      handler.NewTenantTypeHandler(tenantTypeService, userLogService),
		AttractionPriceHandler: handler.NewAttractionPriceHandler(attractionPriceService, userLogService),
		ParkingHandler:         handler.NewParkingHandler(parkingService, userLogService),
		KiosLocationHandler:    handler.NewKiosLocationHandler(kiosLocationService, userLogService),
		TicketingHandler:       handler.NewTicketingHandler(ticketingService, userLogService),
		AttractionHandler:      handler.NewAttractionHandler(attractionService, userLogService),
		TenantHandler:          handler.NewTenantHandler(tenantService, tenantTypeService, kiosLocationService, userLogService),
		KiosAvailableHandler:   handler.NewKiosAvailableHandler(kiosAvailableService),
		BillingHistoryHandler:  handler.NewBillingHistoryHandler(billingService),
		ReportHandler:          handler.NewReportHandler(reportService),
	}

	r := router.New(cfg, h, userLogService)

	log.Printf("server berjalan di port %s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("gagal menjalankan server: %v", err)
	}
}
