package config

import "github.com/BurntSushi/toml"

type Config struct {
	Db   DB   `toml:"db"`
	Coin Coin `toml:"coin"`
}

type Coin struct {
	Id    int    `toml:"id"`
	Name  string `toml:"name"`
	Limit int    `toml:"limit"`
}

type DB struct {
	Addr     string `toml:"addr"`
	Port     string `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Database string `toml:"database"`
}

func New(config *Config, configPath string) error {
	_, err := toml.DecodeFile(configPath, config)
	return err
}
