package router

import (
	"be-capstone-project/src/cmd/public/controller"
	"be-capstone-project/src/internal/core/logger"
	"github.com/gin-gonic/gin"
)

// RegisterGinRouters All router will register here
func RegisterGinRouters(
	in *gin.Engine,
	sampleController *controller.SampleController,
) {

	sampleGroup := in.Group("/sample")
	{
		sampleGroup.POST("", sampleController.CreateSampleController)
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
