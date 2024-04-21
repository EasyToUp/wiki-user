package global

import (
	"github.com/redis/go-redis/v9"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
	"sync"
	"wiki-user/server/config"
	"wiki-user/server/util/timer"
)

var (
	WK_VIPER               *viper.Viper
	WK_LOG                 *zap.Logger
	WK_DBList              map[string]*gorm.DB
	WK_DB                  *gorm.DB
	lock                   sync.RWMutex
	WK_CONFIG              config.Server
	WK_Concurrency_Control = &singleflight.Group{}
	BlackCache             local_cache.Cache
	WK_REDIS               *redis.Client
	WK_Timer               timer.Timer = timer.NewTimerTask()
)

// GetGlobalDBByDBName 通过名称获取db list中的db
func GetGlobalDBByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	return WK_DBList[dbname]
}

// MustGetGlobalDBByDBName 通过名称获取db 如果不存在则panic
func MustGetGlobalDBByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	db, ok := WK_DBList[dbname]
	if !ok || db == nil {
		panic("db no init")
	}
	return db
}
