package routes

import (
	"net/http"
	"time"

	"github.com/Slightly-Techie/st-okr-api/internal/middleware"
	"github.com/Slightly-Techie/st-okr-api/provider"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(prov *provider.Provider) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(ErrorHandlerMiddleware())
	router.Use(gin.Recovery())

	v1 := router.Group("/api/v1")

	// Auth routes
	authRoutes := v1.Group("/auth")
	{
		authRoutes.GET("/:provider", prov.UserController.ContinueWithOAuth)
		authRoutes.GET("/:provider/callback", prov.UserController.GetOAuthCallback)
		authRoutes.GET("/logout/:provider", prov.UserController.LogoutWithOAuth)
	}

	// Company routes
	companyRoutes := v1.Group("/companies")
	companyRoutes.Use(middleware.RequireAuth(prov))
	{
		companyRoutes.POST("/", prov.CompanyController.CreateCompany)
		companyRoutes.GET("/:id", prov.CompanyController.GetCompany)
		companyRoutes.PUT("/:id", prov.CompanyController.UpdateCompany)
		companyRoutes.DELETE("/:id", prov.CompanyController.DeleteCompany)
	}

	// Membership routes
	membershipRoutes := v1.Group("/memberships")
	membershipRoutes.Use(middleware.RequireAuth(prov))
	{
		membershipRoutes.POST("/", prov.MembershipController.CreateMembership)
		membershipRoutes.GET("/:id", prov.MembershipController.GetMembership)
		membershipRoutes.PUT("/:id", prov.MembershipController.UpdateMembership)
		membershipRoutes.DELETE("/:id", prov.MembershipController.DeleteMembership)
		
		// Additional membership routes
		membershipRoutes.GET("/company/:company_id", prov.MembershipController.GetCompanyMembers)
		membershipRoutes.PATCH("/:id/role", prov.MembershipController.UpdateMembershipRole)
		membershipRoutes.PATCH("/:id/status", prov.MembershipController.UpdateMembershipStatus)
	}

	return router
}

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}
}