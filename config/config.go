package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Task    []Task
}

type Task struct {
	Name    string
	Command string
	Args    string
	Dir     string
}

func Init() (*Config, error) {
	var conf Config

	data, err := os.ReadFile("./config.toml")
	if err != nil {
		return &conf, err
	}

	_, err = toml.Decode(string(data), &conf)
	return &conf, err
}
