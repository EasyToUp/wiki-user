package router

import "wiki-user/server/router/user"

type RouterGroup struct {
	User user.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
