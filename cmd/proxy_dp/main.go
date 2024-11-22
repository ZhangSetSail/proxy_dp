package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"rbd_proxy_dp/config"
	"rbd_proxy_dp/pkg/component"
	"rbd_proxy_dp/pkg/server"
	"rbd_proxy_dp/proxy/route"
)

func main() {
	config.Default().
		SetServerName("proxy").
		SetPort("8080").
		SetProxyFlags().
		Parse()
	err := server.New(context.Background()).
		Registry(component.NewLog()).
		Registry(component.NewDB()).
		Registry(component.NewAPI(route.RegisterRoute)).
		Start()
	if err != nil {
		logrus.Errorf("start rbd-api error %s", err.Error())
	}
}
