package server

import (
	"github.com/EAHITechnology/inf/golang/log"
	"github.com/EAHITechnology/inf/unique-number/utils"
)

type Config struct {
	Step               int64 `toml:"step"`
	ExpansionRemaining int64 `toml:"expansion_remaining"`
	LoadFrequency      int   `toml:"load_frequency"`
}

func initConfig(config Config) error {
	//TODO::这里可以进行init相关的数据
	utils.Step = config.Step
	utils.LoadFrequency = config.LoadFrequency
	utils.ExpansionRemaining = config.ExpansionRemaining
	log.Infof("init Config:::%+v", config)
	return nil
}
