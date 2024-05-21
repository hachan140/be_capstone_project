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
		organizationGroup.DELETE("/:id")
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
