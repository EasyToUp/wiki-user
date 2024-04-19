package user

import "wiki-user/server/service"

type ApiGroup struct {
	AccountApi
}

var (
	accountService = service.ServiceGroupApp.AccountServiceGroup.AccountService
)
