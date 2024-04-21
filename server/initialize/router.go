package initialize

import (
	"github.com/gin-gonic/gin"
	"wiki-user/server/global"
	"wiki-user/server/middleware"
	"wiki-user/server/router"
)

func Routers() *gin.Engine {
	Router := gin.New()

	// 跨域，如需跨域可以打开下面的注释
	Router.Use(middleware.Cors()) // 直接放行全部跨域请求
	//Router.Use(middleware.CorsByRules()) // 按照配置的规则放行跨域请求
	//global.WK_LOG.Info("use middleware cors")

	Router.Use(middleware.DefaultLogger())
	publicRouter := Router.Group(global.WK_CONFIG.System.RouterPrefix)
	privateRouter := Router.Group(global.WK_CONFIG.System.RouterPrefix)
	privateRouter.Use(middleware.JWTAuth())

	accountRouter := router.RouterGroupApp.System.AccountRouter
	baseRouter := router.RouterGroupApp.System.BaseRouter
	jwtRouter := router.RouterGroupApp.System.JwtRouter

	{
		accountRouter.InitAccountRouter(publicRouter)
		baseRouter.InitBaseRouter(publicRouter)
		jwtRouter.InitJwtRouter(publicRouter)
	}
	{

	}

	return Router
}
