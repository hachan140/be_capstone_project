package router

import (
	"be-capstone-project/src/cmd/public/controller"
	"be-capstone-project/src/cmd/public/middleware"
	"be-capstone-project/src/internal/core/logger"
	"github.com/gin-gonic/gin"
	"os"
)

// RegisterGinRouters All router will register here
func RegisterGinRouters(
	in *gin.Engine,
	sampleController *controller.SampleController,
	authController *controller.AuthController,
	organizationController *controller.OrganizationController,
	categoryController *controller.CategoryController,
) {
	publicKey := os.Getenv("ACCESS_TOKEN_PUBLIC_KEY")

	sampleGroup := in.Group("/sample")
	sampleGroup.Use(middleware.ValidateToken(publicKey))
	{
		sampleGroup.POST("", sampleController.CreateSampleController)
	}

	authGroup := in.Group("/auth")
	{
		authGroup.POST("/login", authController.Login)
		authGroup.POST("/signup", authController.Signup)
		authGroup.POST("/social-login", authController.SocialLogin)
		authGroup.POST("/refresh-token", authController.RefreshToken)
	}

	organizationGroup := in.Group("/organization")
	organizationGroup.Use(middleware.ValidateToken(publicKey))
	{
		organizationGroup.POST("", organizationController.CreateOrganization)
		organizationGroup.GET("/:id", organizationController.ViewOrganization)
		organizationGroup.PATCH("/:id", organizationController.UpdateOrganization)
		organizationGroup.POST("/:id/add-people", organizationController.AddPeopleToOrganization)
	}

	categoryGroup := in.Group("/category")
	categoryGroup.Use(middleware.ValidateToken(publicKey))
	{
		categoryGroup.POST("", categoryController.CreateCategory)
		categoryGroup.GET("/:id", categoryController.ViewCategoryByID)
		categoryGroup.PATCH("/:id", categoryController.UpdateCategory)
		categoryGroup.GET("/organization/:id", categoryController.ViewListCategoryByOrganization)
	}
	group := in.Group("/test")
	{
		group.GET("", func(context *gin.Context) {
			logger.InfoCtx(context.Request.Context(), "hehe")
			logger.WarnCtx(context, "something warning you")
			logger.ErrorCtx(context, "Error ca ngay")
			context.JSON(200, "hehe")
			return
		})
	}
}
