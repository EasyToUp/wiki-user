package system

import (
	"github.com/gin-gonic/gin"
	"wiki-user/server/api"
)

type AccountRouter struct{}

func (s *AccountRouter) InitAccountRouter(Router *gin.RouterGroup) {
	apiRouter := Router.Group("account")
	accountApi := api.ApiGroupApp.SystemApiGroup.AccountApi
	{
		apiRouter.POST("login", accountApi.Login)       // login
		apiRouter.POST("register", accountApi.Register) // register
	}
}
