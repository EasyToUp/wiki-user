package config

type System struct {
	Addr         int    `mapstructure:"addr" json:"addr" yaml:"addr"`       // 端口值
	ApiUrl       string `mapstructure:"apiUrl" json:"apiUrl" yaml:"apiurl"` // wiki-api-url
	RouterPrefix string `mapstructure:"router-prefix" json:"router-prefix" yaml:"router-prefix"`
}
