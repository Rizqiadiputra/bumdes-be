package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/liyansasongko/bumdes-be/internal/config"
	"github.com/liyansasongko/bumdes-be/internal/handler"
	"github.com/liyansasongko/bumdes-be/internal/middleware"
	"github.com/liyansasongko/bumdes-be/internal/service"
)

type Handlers struct {
	AuthHandler            *handler.AuthHandler
	UserHandler            *handler.UserHandler
	RoleHandler            *handler.RoleHandler
	PermissionHandler      *handler.PermissionHandler
	UserLogHandler         *handler.UserLogHandler
	ParkingPriceHandler    *handler.ParkingPriceHandler
	PaymentMethodHandler   *handler.PaymentMethodHandler
	RevenueCategoryHandler *handler.RevenueCategoryHandler
	TicketPriceHandler     *handler.TicketPriceHandler
	TenantTypeHandler      *handler.TenantTypeHandler
	AttractionPriceHandler *handler.AttractionPriceHandler
	ParkingHandler         *handler.ParkingHandler
	KiosLocationHandler    *handler.KiosLocationHandler
	TicketingHandler       *handler.TicketingHandler
	AttractionHandler      *handler.AttractionHandler
	TenantHandler          *handler.TenantHandler
	KiosAvailableHandler   *handler.KiosAvailableHandler
	BillingHistoryHandler  *handler.BillingHistoryHandler
	ReportHandler          *handler.ReportHandler
}

func New(cfg *config.Config, h Handlers, userLogService service.UserLogService) *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", h.AuthHandler.Login)
		}

		protected := api.Group("")
		protected.Use(middleware.Auth(cfg.JWTSecret))
		protected.Use(middleware.ActivityLog(userLogService))
		{
			protected.GET("/me", h.UserHandler.Me)

			roles := protected.Group("/roles")
			{
				roles.GET("", h.RoleHandler.ListRoles)
				roles.GET("/:id", h.RoleHandler.GetRole)
				roles.PUT("/:id/permissions", h.RoleHandler.UpdatePermissions)
			}

			permissions := protected.Group("/permissions")
			{
				permissions.GET("", h.PermissionHandler.ListPermissions)
			}

			accounts := protected.Group("/accounts")
			{
				accounts.GET("", h.UserHandler.ListAccounts)
				accounts.POST("", h.UserHandler.CreateAccount)
				accounts.PUT("/:id", h.UserHandler.UpdateAccount)
				accounts.DELETE("/:id", h.UserHandler.DeleteAccount)
			}

			logs := protected.Group("/logs")
			{
				logs.GET("", h.UserLogHandler.ListLogs)
			}

			masterData := protected.Group("/master-data")
			{
				parkingPrices := masterData.Group("/parking-prices")
				{
					parkingPrices.GET("", h.ParkingPriceHandler.ListParkingPrices)
					parkingPrices.GET("/:id", h.ParkingPriceHandler.GetParkingPrice)
					parkingPrices.POST("", h.ParkingPriceHandler.CreateParkingPrice)
					parkingPrices.PUT("/:id", h.ParkingPriceHandler.UpdateParkingPrice)
					parkingPrices.DELETE("/:id", h.ParkingPriceHandler.DeleteParkingPrice)
				}

				paymentMethods := masterData.Group("/payment-methods")
				{
					paymentMethods.GET("", h.PaymentMethodHandler.ListPaymentMethods)
					paymentMethods.GET("/:id", h.PaymentMethodHandler.GetPaymentMethod)
					paymentMethods.POST("", h.PaymentMethodHandler.CreatePaymentMethod)
					paymentMethods.PUT("/:id", h.PaymentMethodHandler.UpdatePaymentMethod)
					paymentMethods.DELETE("/:id", h.PaymentMethodHandler.DeletePaymentMethod)
				}

				revenueCategories := masterData.Group("/revenue-categories")
				{
					revenueCategories.GET("", h.RevenueCategoryHandler.ListRevenueCategories)
					revenueCategories.GET("/:id", h.RevenueCategoryHandler.GetRevenueCategory)
					revenueCategories.POST("", h.RevenueCategoryHandler.CreateRevenueCategory)
					revenueCategories.PUT("/:id", h.RevenueCategoryHandler.UpdateRevenueCategory)
					revenueCategories.DELETE("/:id", h.RevenueCategoryHandler.DeleteRevenueCategory)
				}

				ticketPrices := masterData.Group("/ticket-prices")
				{
					ticketPrices.GET("", h.TicketPriceHandler.ListTicketPrices)
					ticketPrices.GET("/:id", h.TicketPriceHandler.GetTicketPrice)
					ticketPrices.POST("", h.TicketPriceHandler.CreateTicketPrice)
					ticketPrices.PUT("/:id", h.TicketPriceHandler.UpdateTicketPrice)
					ticketPrices.DELETE("/:id", h.TicketPriceHandler.DeleteTicketPrice)
				}

				tenantTypes := masterData.Group("/tenant-types")
				{
					tenantTypes.GET("", h.TenantTypeHandler.ListTenantTypes)
					tenantTypes.GET("/:id", h.TenantTypeHandler.GetTenantType)
					tenantTypes.POST("", h.TenantTypeHandler.CreateTenantType)
					tenantTypes.PUT("/:id", h.TenantTypeHandler.UpdateTenantType)
					tenantTypes.DELETE("/:id", h.TenantTypeHandler.DeleteTenantType)
				}

				attractionPrices := masterData.Group("/attraction-prices")
				{
					attractionPrices.GET("", h.AttractionPriceHandler.ListAttractionPrices)
					attractionPrices.GET("/:id", h.AttractionPriceHandler.GetAttractionPrice)
					attractionPrices.POST("", h.AttractionPriceHandler.CreateAttractionPrice)
					attractionPrices.PUT("/:id", h.AttractionPriceHandler.UpdateAttractionPrice)
					attractionPrices.DELETE("/:id", h.AttractionPriceHandler.DeleteAttractionPrice)
				}

				kiosLocations := masterData.Group("/kios_locations")
				{
					kiosLocations.GET("", h.KiosLocationHandler.ListKiosLocations)
					kiosLocations.GET("/:id", h.KiosLocationHandler.GetKiosLocation)
					kiosLocations.POST("", h.KiosLocationHandler.CreateKiosLocation)
					kiosLocations.PUT("/:id", h.KiosLocationHandler.UpdateKiosLocation)
					kiosLocations.DELETE("/:id", h.KiosLocationHandler.DeleteKiosLocation)
				}
			}

			ticketing := protected.Group("/ticketing")
			{
				ticketing.GET("", h.TicketingHandler.ListTicketings)
				ticketing.GET("/:id", h.TicketingHandler.GetTicketing)
				ticketing.POST("", h.TicketingHandler.CreateTicketing)
				ticketing.PUT("/:id", h.TicketingHandler.UpdateTicketing)
				ticketing.DELETE("/:id", h.TicketingHandler.DeleteTicketing)
			}

			parkings := protected.Group("/parkings")
			{
				parkings.GET("", h.ParkingHandler.ListParkings)
				parkings.GET("/:id", h.ParkingHandler.GetParking)
				parkings.POST("", h.ParkingHandler.CreateParking)
				parkings.DELETE("/:id", h.ParkingHandler.DeleteParking)
			}

			attractions := protected.Group("/attractions")
			{
				attractions.GET("", h.AttractionHandler.ListAttractions)
				attractions.GET("/:id", h.AttractionHandler.GetAttraction)
				attractions.POST("", h.AttractionHandler.CreateAttraction)
				attractions.PUT("/:id", h.AttractionHandler.UpdateAttraction)
				attractions.DELETE("/:id", h.AttractionHandler.DeleteAttraction)
			}

			tenants := protected.Group("/tenants")
			{
				tenants.GET("", h.TenantHandler.ListTenants)
				tenants.GET("/tenant-types", h.TenantHandler.SearchTenantTypes)
				tenants.GET("/kios_locations", h.TenantHandler.SearchKiosLocations)
				tenants.GET("/:id", h.TenantHandler.GetTenant)
				tenants.POST("", h.TenantHandler.CreateTenant)
				tenants.PUT("/:id", h.TenantHandler.UpdateTenant)
				tenants.DELETE("/:id", h.TenantHandler.DeleteTenant)
			}

			kiosAvailables := protected.Group("/kios_availables")
			{
				kiosAvailables.GET("", h.KiosAvailableHandler.ListKiosAvailables)
			}

			billingHistories := protected.Group("/billing-histories")
			{
				billingHistories.GET("", h.BillingHistoryHandler.ListBillingHistories)
				billingHistories.GET("/:id", h.BillingHistoryHandler.GetBillingHistory)
			}

			reports := protected.Group("/reports")
			{
				reports.GET("/revenue", h.ReportHandler.GetRevenueReport)
			}
		}
	}

	return r
}
