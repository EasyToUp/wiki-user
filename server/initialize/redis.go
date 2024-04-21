package initialize

import (
	"context"

	"wiki-user/server/global"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func Redis() {
	redisCfg := global.WK_CONFIG.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password, // no password set
		DB:       redisCfg.DB,       // use default DB
	})
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		global.WK_LOG.Error("redis connect ping failed, err:", zap.Error(err))
	} else {
		global.WK_LOG.Info("redis connect ping response:", zap.String("pong", pong))
		global.WK_REDIS = client
	}
}
