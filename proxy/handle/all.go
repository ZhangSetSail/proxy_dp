package handle

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"net/url"
	"rbd_proxy_dp/config"
)

func ProxyRouteHandle(w http.ResponseWriter, r *http.Request) {
	proxyTarget := config.DefaultProxy().ProxyTarget

	// 解析代理目标 URL
	proxyURL, err := url.Parse(proxyTarget)
	if err != nil {
		http.Error(w, `{"error_msg":"invalid proxy URL"}`, http.StatusInternalServerError)
		logrus.Errorf("error parsing proxy URL: %v", err)
		return
	}

	// 创建反向代理
	proxy := httputil.NewSingleHostReverseProxy(proxyURL)

	// 自定义代理请求处理
	proxy.ModifyResponse = func(resp *http.Response) error {
		logrus.Infof("response from target: status=%d, url=%s", resp.StatusCode, resp.Request.URL)
		return nil
	}
	proxy.ErrorHandler = func(w http.ResponseWriter, req *http.Request, err error) {
		logrus.Errorf("proxy error: %v", err)
		http.Error(w, `{"error_msg":"proxy error"}`, http.StatusBadGateway)
	}
	proxy.Director = func(req *http.Request) {
		// 设置目标地址
		req.URL.Scheme = proxyURL.Scheme
		req.URL.Host = proxyURL.Host
		req.Host = proxyURL.Host // 更新 Host Header
		logrus.Infof("forwarding request: method=%s, path=%s, target=%s", req.Method, req.URL.Path, proxyURL)
	}
	// 执行代理
	proxy.ServeHTTP(w, r)
}
