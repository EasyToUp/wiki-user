package initialize

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
	"wiki-user/server/global"
	"wiki-user/server/model/example"
)

func Gorm() *gorm.DB {
	switch global.WK_CONFIG.System.DbType {
	case "mysql":
		return GormMysql()
	case "pgsql":
		return GormPgSql()
	case "oracle":
		return GormOracle()
	case "mssql":
		return GormMssql()
	case "sqlite":
		return GormSqlite()
	default:
		return GormMysql()
	}
}

func RegisterTables() {
	db := global.WK_DB
	err := db.AutoMigrate(
		example.ExaFileUploadAndDownload{},
	)
	if err != nil {
		global.WK_LOG.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}
	global.WK_LOG.Info("register table success")
}
