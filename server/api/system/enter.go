package system

import "wiki-user/server/service"

type ApiGroup struct {
	JwtApi
	BaseApi
	AccountApi
}

var (
	jwtService     = service.ServiceGroupApp.SystemServiceGroup.JwtService
	accountService = service.ServiceGroupApp.AccountServiceGroup.AccountService
)
