package conf

import (
	"github.com/letterScape/backend/global"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	vp *viper.Viper
}

func NewConfig() (*Config, error) {
	vp := viper.New()
	vp.AddConfigPath("conf/dev")
	vp.SetConfigName("config")
	vp.SetConfigType("yaml")
	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &Config{vp: vp}, nil
}

func (c *Config) ReadSection(k string, v interface{}) error {
	err := c.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}
	return nil
}

func SetupConfig() {
	config, err := NewConfig()
	if err != nil {
		log.Panic("NewConfig error", err)
	}
	err = config.ReadSection("BlockChain", &global.BlockChainConfig)
	if err != nil {
		log.Panic("ReadSection-BlockChain error: ", err)
	}

}
