package api

import (
	"wiki-user/server/api/system"
	"wiki-user/server/api/user"
)

type ApiGroup struct {
	AccountApiGroup user.ApiGroup
	SystemApiGroup  system.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
