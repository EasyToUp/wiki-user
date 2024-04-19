package main

import (
	"wiki-user/server/core"
	"wiki-user/server/global"
)

func main() {
	global.VIPER = core.LoadConfig()
	core.RunServer()
}
