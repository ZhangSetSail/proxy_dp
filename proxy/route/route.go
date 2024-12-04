package route

import (
	"net/http"
	"rbd_proxy_dp/proxy/handle"
)

func RegisterRoute(mux *http.ServeMux) {
	// 未匹配路径走代理
	mux.HandleFunc("/openapi/v1/monitor/proxy_domain", handle.ProxyDomain)
	mux.HandleFunc("/", handle.ProxyRouteHandle)
}
