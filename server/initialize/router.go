package initialize

import (
	"github.com/gin-gonic/gin"
	"wiki-user/server/global"
	"wiki-user/server/middleware"
	"wiki-user/server/router"
)

func Routers() *gin.Engine {
	Router := gin.New()
	Router.Use(middleware.DefaultLogger())
	baseGroup := Router.Group(global.WK_CONFIG.System.RouterPrefix)
	accountRouter := router.RouterGroupApp.User
	accountRouter.InitAccountRouter(baseGroup)
	return Router
}
