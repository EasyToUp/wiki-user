package main

import (
	"go.uber.org/zap"
	"wiki-user/server/core"
	"wiki-user/server/global"
	"wiki-user/server/initialize"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

func main() {
	global.WK_VIPER = core.Viper()
	global.WK_LOG = core.Zap() // 初始化zap日志库
	zap.ReplaceGlobals(global.WK_LOG)
	global.WK_DB = initialize.Gorm() // gorm连接数据库
	initialize.Timer()
	initialize.DBList()
	if global.WK_DB != nil {
		initialize.RegisterTables() // 初始化表
		initialize.CheckSetting()   // 检测设置
		// 程序结束前关闭数据库链接
		db, _ := global.WK_DB.DB()
		defer db.Close()
	}
	core.RunServer()
}
