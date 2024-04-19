package core

import (
	"fmt"
	"wiki-user/server/global"
	"wiki-user/server/initialize"
)

func RunServer() {

	Router := initialize.Routers()
	address := fmt.Sprintf(":%d", global.CONFIG.System.Addr)
	Router.Run(address)
	fmt.Println("server run success on ", address)

}
