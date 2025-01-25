package config

import (
	"log"
	"os"
	"slices"

	"gopkg.in/yaml.v3"
)

var ALGORITHMS = []string{"round_robin"}

func LoadSystemConfig() SystemEnv {
	configByte, err := os.ReadFile("./.env.yaml")
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}
	var config SystemEnv
	err = yaml.Unmarshal(configByte, &config)
	if err != nil {
		log.Fatalf("failed to unmarshal config file: %v", err)
	}
	return config
}

func LoadProxyConfig() ProxyMapping {
	configByte, err := os.ReadFile("./proxy-setting.yaml")
	if err != nil {
		log.Fatalf("failed to read proxy config file: %v", err)
	}
	var config ProxyMapping
	// add a check to validate if the host names are all different
	err = yaml.Unmarshal(configByte, &config)
	if err != nil {
		log.Fatalf("failed to unmarshal proxy config file: %v", err)
	}
	isAlgoValid := slices.Contains(ALGORITHMS, config.Algorithm)
	if !isAlgoValid {
		log.Fatalf("invalid algorithm: %v", config.Algorithm)
	}
	return config
}
