package service

import "wiki-user/server/service/user"

type ServiceGroup struct {
	AccountServiceGroup user.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
