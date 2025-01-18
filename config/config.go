package config

import (
	"os"
	"reverse_proxy/utils"
	"slices"

	"gopkg.in/yaml.v3"
)

func LoadSystemConfig() SystemEnv {
	configByte, err := os.ReadFile("./.env.yaml")
	if err != nil {
		utils.JsonLog("failed to read config file", "fatal", err)
	}
	var config SystemEnv
	err = yaml.Unmarshal(configByte, &config)
	if err != nil {
		utils.JsonLog("failed to unmarshal config file", "fatal", err)
	}
	return config
}

func LoadProxyConfig() ProxyMapping {
	configByte, err := os.ReadFile("./proxy-setting.yaml")
	if err != nil {
		utils.JsonLog("failed to read proxy config file", "fatal", err)
	}
	var config ProxyMapping
	err = yaml.Unmarshal(configByte, &config)
	if err != nil {
		utils.JsonLog("failed to unmarshal proxy config file", "fatal", err)
	}
	isAlgoValid := slices.Contains(utils.ALGORITHMS, config.Algorithm)
	if !isAlgoValid {
		utils.JsonLog("invalid algorithm", "fatal", nil)
	}
	return config
}
