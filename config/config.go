package config

import (
	"rbd_proxy_dp/config/public"
	"rbd_proxy_dp/config/server"
)

func init() {
	defaultConfig = &Config{
		ProxyConfig: &server.ProxyConfig{
			ProxyTarget: "https://large-screen.goodrain.com",
		},
		PublicConfig: &public.PublicConfig{
			LogLevel:   "info",
			Port:       "",
			ServerName: "",
		},
		DBConfig: &public.DBConfig{
			DBType:   "sqlite",
			DBHost:   "",
			DBPort:   "",
			DBUser:   "",
			DBPasswd: "",
			DBName:   "proxy",
		},
	}
}

// Config -
type Config struct {
	ProxyConfig  *server.ProxyConfig
	PublicConfig *public.PublicConfig
	DBConfig     *public.DBConfig
}

var defaultConfig *Config

// Default -
func Default() *Config {
	return defaultConfig
}

// DefaultProxy -
func DefaultProxy() *server.ProxyConfig {
	return defaultConfig.ProxyConfig
}

// DefaultPublic -
func DefaultPublic() *public.PublicConfig {
	return defaultConfig.PublicConfig
}

func DefaultDB() *public.DBConfig {
	return defaultConfig.DBConfig
}
