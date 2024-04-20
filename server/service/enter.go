package service

import (
	"wiki-user/server/service/system"
	"wiki-user/server/service/user"
)

type ServiceGroup struct {
	AccountServiceGroup user.ServiceGroup
	SystemServiceGroup  system.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
