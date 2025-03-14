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
	hyperDocumentController *controller.HyperDocumentController,
	productController *controller.ProductController,
) {
	publicKey := os.Getenv("ACCESS_TOKEN_PUBLIC_KEY")
	publicKeyResetPassword := os.Getenv("RESET_PASSWORD_PUBLIC_KEY")

	sampleGroup := in.Group("/sample")
	sampleGroup.Use(middleware.ValidateToken(publicKey))
	{
		sampleGroup.POST("", sampleController.CreateSampleController)
	}

	authGroup := in.Group("/auth")
	{
		authGroup.POST("/login", authController.Login)
		authGroup.POST("/signup", authController.Signup)
		authGroup.GET("/email/verify", authController.VerifyEmail)
		authGroup.POST("/social-login", authController.SocialLogin)
		authGroup.POST("/refresh-token", authController.RefreshToken)
		authGroup.POST("/reset-password/request", authController.ResetPasswordRequest)
	}
	resetPassword := in.Group("/reset-password")
	resetPassword.Use(middleware.ValidateToken(publicKeyResetPassword))
	{
		resetPassword.POST("", authController.ResetPassword)
	}

	organizationGroup := in.Group("/organization")
	organizationGroup.Use(middleware.ValidateToken(publicKey))
	{
		organizationGroup.POST("", organizationController.CreateOrganization)
		organizationGroup.GET("/:id", organizationController.ViewOrganization)
		organizationGroup.PATCH("/:id", organizationController.UpdateOrganization)
		organizationGroup.PATCH("/:id/status", organizationController.UpdateOrganizationStatus)
		organizationGroup.POST("/:id/add-people", organizationController.AddPeopleToOrganization)
		organizationGroup.POST(":id/remove-people", organizationController.RemovePeopleFromOrganization)
		organizationGroup.POST("/:id/assign-manager", organizationController.AssignPeopleToManager)
		organizationGroup.POST("/:id/recall-manager", organizationController.RecallPeopleToManager)
		organizationGroup.GET("/check-user-organization", organizationController.CheckUserAlreadyRequestCreateOrganization)
	}
	acceptInvitation := in.Group("/accept")
	{
		acceptInvitation.GET("/:orgID/dept/:deptID/user/:userEmail", organizationController.AcceptOrganizationInvitation)
	}

	categoryGroup := in.Group("/category")
	categoryGroup.Use(middleware.ValidateToken(publicKey))
	{
		categoryGroup.POST("", categoryController.CreateCategory)
		categoryGroup.GET("/:id", categoryController.ViewCategoryByID)
		categoryGroup.PATCH("/:id", categoryController.UpdateCategory)
		categoryGroup.PATCH("/:id/status", categoryController.UpdateCategoryStatus)
	}

	departmentGroup := in.Group("/department")
	departmentGroup.Use(middleware.ValidateToken(publicKey))
	{
		departmentGroup.GET("/:id/category/by-name", categoryController.ViewCategoryByNameLike)
		departmentGroup.GET("/:id/category", categoryController.ViewListCategoryByDepartment)
		departmentGroup.PATCH("/:id/status", categoryController.UpdateDepartmentStatus)
	}
	documentGroup := in.Group("/document")
	documentGroup.Use(middleware.ValidateToken(publicKey))
	{
		documentGroup.GET("", hyperDocumentController.FilterHyperDocument)
		documentGroup.POST("/search/and-or-not", hyperDocumentController.SearchDocumentAndOrNot)
		documentGroup.POST("/search/keyword", hyperDocumentController.GetSearchHistoryKeywords)
		documentGroup.POST("/search-history", hyperDocumentController.SaveSearchHistory)
	}

	productGroup := in.Group("/product")
	productGroup.Use(middleware.ValidateToken(publicKey))
	{
		productGroup.POST("", productController.CreateProduct)
		productGroup.GET("", productController.ViewListProduct)
		productGroup.GET("/:id", productController.GetProductByID)
		productGroup.GET("/by-name/:name", productController.SearchProductByName)
		productGroup.PATCH("/:id", productController.UpdateProduct)
		productGroup.DELETE("/:id", productController.DeleteProduct)
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
