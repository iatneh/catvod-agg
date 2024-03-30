package conf

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	Logger  *LogConfig
	General *GeneralConfig
}

// LogConfig 日志配置
type LogConfig struct {
	LogLevel string // 配置日志输出级别: trace,debug,info,warn,error
}

// GeneralConfig 通用配置
type GeneralConfig struct {
	FileName   string // 多站配置文件
	ToFilePath string // 聚合结果输出路径
	OtherSites string // 其他单仓地址
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

	if AppConf.Logger != nil && len(AppConf.Logger.LogLevel) == 0 {
		AppConf.Logger.LogLevel = "debug"
	}
	ll, err := logrus.ParseLevel(AppConf.Logger.LogLevel)
	if err != nil {
		ll = logrus.DebugLevel
	}
	logrus.SetLevel(ll)
}
