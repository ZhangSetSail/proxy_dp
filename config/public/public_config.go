package public

import (
	"github.com/alecthomas/kingpin/v2"
)

// PublicConfig config
type PublicConfig struct {
	LogLevel   string
	Port       string
	ServerName string
}

func AddPublicFlags(puc *PublicConfig) {
	kingpin.Flag("log-level", "The level of logger").Default(puc.LogLevel).Envar("LOG_LEVEL").StringVar(&puc.LogLevel)
	kingpin.Flag("port", "port").Default(puc.Port).Envar("PORT").StringVar(&puc.Port)
	kingpin.Flag("server-name", "server name").Default(puc.ServerName).Envar("SERVER_NAME").StringVar(&puc.ServerName)
}
