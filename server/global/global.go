package global

import (
	"github.com/go-redis/redis/v8"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
	"sync"
	"wiki-user/server/config"
)

var (
	CONFIG                  config.Server
	VIPER                   *viper.Viper
	GVA_LOG                 *zap.Logger
	GVA_DBList              map[string]*gorm.DB
	GVA_DB                  *gorm.DB
	lock                    sync.RWMutex
	GVA_CONFIG              config.Server
	GVA_Concurrency_Control = &singleflight.Group{}
	BlackCache              local_cache.Cache
	GVA_REDIS               *redis.Client
)

// GetGlobalDBByDBName 通过名称获取db list中的db
func GetGlobalDBByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	return GVA_DBList[dbname]
}

// MustGetGlobalDBByDBName 通过名称获取db 如果不存在则panic
func MustGetGlobalDBByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	db, ok := GVA_DBList[dbname]
	if !ok || db == nil {
		panic("db no init")
	}
	return db
}
