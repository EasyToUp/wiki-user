package user

import (
	"github.com/gin-gonic/gin"
	"wiki-user/server/api"
	"wiki-user/server/middleware"
)

type AccountRouter struct{}

func (s *AccountRouter) InitAccountRouter(Router *gin.RouterGroup) {
	apiRouter := Router.Group("account")
	apiRouter.Use(middleware.CORSMiddleware())
	apiRouter.Use(middleware.SessionManager())
	accountApi := api.ApiGroupApp.AccountApiGroup.AccountApi
	{
		apiRouter.POST("login", accountApi.Login)       // 创建Api
		apiRouter.POST("register", accountApi.Register) // 删除Api

	}
}
