package initialize

import (
	"github.com/gin-gonic/gin"
	"wiki-user/server/global"
	"wiki-user/server/router"
)

func Routers() *gin.Engine {
	Router := gin.New()
	baseGroup := Router.Group(global.CONFIG.System.RouterPrefix)
	accountRouter := router.RouterGroupApp.User
	accountRouter.InitAccountRouter(baseGroup)
	return Router
}
