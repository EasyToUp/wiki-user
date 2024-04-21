package system

import (
	"github.com/gin-gonic/gin"
	"wiki-user/server/api"
	"wiki-user/server/middleware"
)

type AccountRouter struct{}

func (s *AccountRouter) InitAccountRouter(Router *gin.RouterGroup) {
	apiRouter := Router.Group("account")
	apiRouter.Use(middleware.SessionManager()) // todo 优化
	accountApi := api.ApiGroupApp.SystemApiGroup.AccountApi
	{
		apiRouter.POST("login", accountApi.Login)       // login
		apiRouter.POST("register", accountApi.Register) // register
	}
}
