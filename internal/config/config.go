package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Bot BotConfig      `json:"bot"`
	Db  DatabaseConfig `json:"database,db"`
}

type BotConfig struct {
	Token             string `json:"token"`
	Debug             bool   `json:"debug"`
	Timeout           int    `json:"timeout,omitempty"`
	Concurrent        bool   `json:"concurrent,omitempty"`
	UpstreamChannelId string `json:"upstream_channel"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func NewConfig(filepath string) (*Config, error) {
	byteData, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(byteData, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
