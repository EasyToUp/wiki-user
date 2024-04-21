package core

import (
	"fmt"
	"wiki-user/server/global"
	"wiki-user/server/initialize"
	"wiki-user/server/service/system"
)

func RunServer() {
	if global.WK_CONFIG.System.UseMultipoint || global.WK_CONFIG.System.UseRedis {
		// 初始化redis服务
		initialize.Redis()
	}
	// 从db加载jwt数据
	if global.WK_DB != nil {
		system.LoadAll()
	}
	Router := initialize.Routers()
	address := fmt.Sprintf(":%d", global.WK_CONFIG.System.Addr)
	Router.Run(address)
	fmt.Println("server run success on ", address)

}
