package router

import (
	"wiki-user/server/router/system"
)

type RouterGroup struct {
	System system.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
