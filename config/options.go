package config

import (
	"github.com/alecthomas/kingpin/v2"
	"rbd_proxy_dp/config/public"
	"rbd_proxy_dp/config/server"
)

func (c *Config) SetProxyFlags() *Config {
	public.AddDBFlags(c.DBConfig)
	server.AddProxyFlags(c.ProxyConfig)
	public.AddPublicFlags(c.PublicConfig)
	return c
}

func (c *Config) SetServerName(serverName string) *Config {
	c.PublicConfig.ServerName = serverName
	return c
}

func (c *Config) SetPort(port string) *Config {
	c.PublicConfig.Port = port
	return c
}

func (c *Config) Parse() {
	kingpin.Parse()
}
