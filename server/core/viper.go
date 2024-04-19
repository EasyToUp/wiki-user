package core

import (
	"fmt"
	"github.com/spf13/viper"
	"wiki-user/server/global"
)

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() *viper.Viper {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	if err = v.Unmarshal(&global.CONFIG); err != nil {
		fmt.Println(err)
	}
	return v

}
