package api

import (
    "github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
    swaggerFiles "github.com/swaggo/files"
	
    "messenger/internal/api/handler"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()
    h := handler.NewHandler()

    api := r.Group("/api")
    {
        api.POST("/start", h.StartDispatcher)
        api.POST("/stop", h.StopDispatcher)
        api.GET("/sent-messages", h.GetSentMessages)
    }

    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    return r
}
