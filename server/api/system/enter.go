package system

import "wiki-user/server/service"

type ApiGroup struct {
	JwtApi
}

var (
	jwtService = service.ServiceGroupApp.SystemServiceGroup.JwtService
)
