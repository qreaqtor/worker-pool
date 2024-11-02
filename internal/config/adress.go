package config

import (
	"fmt"
	"net"
)

type Adress struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (a *Adress) ToString() string {
	return net.JoinHostPort(a.Host, fmt.Sprint(a.Port))
}
