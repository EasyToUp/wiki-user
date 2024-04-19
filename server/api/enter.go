package api

import (
	"wiki-user/server/api/user"
)

type ApiGroup struct {
	AccountApiGroup user.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
