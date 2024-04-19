package global

import (
	"github.com/spf13/viper"
	"wiki-user/server/config"
)

var (
	CONFIG config.Server
	VIPER  *viper.Viper
)
