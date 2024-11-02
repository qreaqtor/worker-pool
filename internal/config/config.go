package config

import "time"

type Config struct {
	Addr Adress `yaml:"addr"`

	AppType string `yaml:"app_type"`

	Workers WorkersConfig `yaml:"workers"`
}

type WorkersConfig struct {
	Sleeper SleepWorkerConfig `yaml:"sleeper"`
}

type SleepWorkerConfig struct {
	Sleep time.Duration `yaml:"sleep"`
}
