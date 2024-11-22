package server

import "github.com/alecthomas/kingpin/v2"

// ProxyConfig config
type ProxyConfig struct {
	ProxyTarget string
}

func AddProxyFlags(prc *ProxyConfig) {
	kingpin.Flag("proxy-target", "proxy target").Default(prc.ProxyTarget).Envar("PROXY-TARGET").StringVar(&prc.ProxyTarget)
}
