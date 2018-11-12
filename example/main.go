// swagger:meta
//
package main

import (
	"github.com/gin-gonic/gin"

	"github.com/wusphinx/gin-swagger/example/from_request"
	"github.com/wusphinx/gin-swagger/example/test"
	"github.com/wusphinx/gin-swagger/example/test3"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/test", test.Auth(), test3.Test3)
	router.GET("/auto", from_request.FromRequest(from_request.GetUser{}))

	userRouter := router.Group("/user", test.AuthMiddleware)

	test.SetupUserRoutes(userRouter)
}

func main() {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.POST("/", test.Test)

	SetupRoutes(router)

	//userRouter := router.Group("/user2")
	//test.SetupUserRoutes(userRouter)

	router.Run(":8080")
}
