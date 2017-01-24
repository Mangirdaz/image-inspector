package main

import "github.com/mangirdaz/image-inspector/pkg/config"

type KeyValueStorageConfig struct {
	Ip   string
	Port string
}

type ErrorMessage struct {
	Code    int
	Message string
}

type image struct {
	Images []struct {
		Name string `json:"name"`
	} `json:"images"`
}

func InitKeyValueStorageConfig() KeyValueStorageConfig {

	var configuration KeyValueStorageConfig
	configuration.Ip = config.Get("EnvKVStorageIp")
	configuration.Port = config.Get("EnvKVStoragePort")

	return configuration
}
