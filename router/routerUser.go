package router2

import "github.com/gin-gonic/gin"

func RouterLoadUser(engine *gin.Engine)  {
	userRouter := engine.Group("/user")
	{
		userRouter.GET("/userinfo")
		userRouter.POST("/userinfo")
		userRouter.PUT("userinfo")
		userRouter.DELETE("userinfo")
	}
}
