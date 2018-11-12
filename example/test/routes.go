package test

import (
	"github.com/gin-gonic/gin"
	"github.com/wusphinx/gin-swagger/example/test2"
)

func SetupUserRoutes(userRouter *gin.RouterGroup) {
	userRouteWith := userRouter.Group("/test")
	userRouteWith2 := userRouteWith.Group("")
	{
		userRouteWith2.GET("/:name/:action", test2.Test2)
	}
}
