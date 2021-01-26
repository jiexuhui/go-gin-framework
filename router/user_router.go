package router

import (
	v1 "livefun/api/v1"

	"github.com/gin-gonic/gin"
)

// InitRouter user router
func InitRouter(Router *gin.RouterGroup) {

	Router.GET("users", v1.Users)
	UserRouter := Router.Group("user")
	{
		UserRouter.DELETE("s", v1.Users)
	}

}
