package router

import (
	v1 "livefun/api/v1"

	"github.com/gin-gonic/gin"
)

// InitBaseRouter is base router
func InitBaseRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	BaseRouter := Router.Group("base")
	{
		BaseRouter.POST("register", v1.Register)
	}
	return BaseRouter
}
