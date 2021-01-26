package initialize

import (
	_ "livefun/docs" // swager docs
	"livefun/global"
	"livefun/middleware"
	"livefun/router"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// Routers  路由注册
func Routers() *gin.Engine {
	var Router = gin.Default()

	global.LF_LOG.Info("use middleware logger")
	// 跨域
	Router.Use(middleware.Cors())
	global.LF_LOG.Info("use middleware cors")
	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	global.LF_LOG.Info("register swagger handler")

	APIGroup := Router.Group("")

	router.InitRouter(APIGroup)
	router.InitBaseRouter(APIGroup)
	return Router
}
