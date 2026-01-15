package main

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	//读取配置
	InitViperWatch()
	server := InitWireServer()
	err := server.Run(":8080")
	if err != nil {
		panic(err)
	}
}

// InitViperWatch 读取配置，可实现热修改
func InitViperWatch() {
	cf := pflag.String("config", "Config/dev.yaml", "配置文件路径")
	// 用于解析参数，这一步以后，变量才有值
	pflag.Parse()
	viper.SetConfigType("yaml")
	viper.SetConfigFile(*cf)
	viper.WatchConfig()
	//读取配置
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
