package conf

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	FileName string
}

var (
	AppConf *Config
)

func init() {
	v := viper.New()
	v.SetConfigName("app")
	v.AddConfigPath("conf")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal errors config file: %s \n", err))
	}
	if err := v.Unmarshal(&AppConf); err != nil {
		panic(fmt.Errorf("Fatal errors config file: %s \n", err))
	}
}
