package conf

import (
	"fmt"

	"github.com/spf13/viper"
)

// 配置项
type Config struct {
	Global GlobalOps
}

// 全局配置
type GlobalOps struct {
	LibvirtConnectUrl string
	Address           string
}

// 初始化配置文件
func (c *Config) Init(path string) error {
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("abnormal reading of configuration file, %e", err)
	}
	// 将配置文件加载到结构体中
	err = viper.Unmarshal(&c)
	if err != nil {
		return err
	}
	return nil
}
