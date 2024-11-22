package handle

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"rbd_proxy_dp/config"
	"rbd_proxy_dp/model"
	"rbd_proxy_dp/pkg/component"
	"time"
)

// ProxyRouteHandle 代理处理请求并存储数据
func ProxyRouteHandle(w http.ResponseWriter, r *http.Request) {
	proxyTarget := config.DefaultProxy().ProxyTarget
	proxyURL, err := url.Parse(proxyTarget)
	if err != nil {
		http.Error(w, `{"error_msg":"invalid proxy URL"}`, http.StatusInternalServerError)
		logrus.Errorf("error parsing proxy URL: %v", err)
		return
	}

	// 创建反向代理
	proxy := httputil.NewSingleHostReverseProxy(proxyURL)

	// 捕获请求参数
	requestBody, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(requestBody)) // 重置请求体供后续读取

	// 自定义代理响应处理
	proxy.ModifyResponse = func(resp *http.Response) error {
		return handleResponse(r, requestBody, resp)
	}

	// 自定义错误处理
	proxy.ErrorHandler = func(w http.ResponseWriter, req *http.Request, err error) {
		handleProxyError(w, req, requestBody, err)
	}

	// 自定义请求处理
	proxy.Director = func(req *http.Request) {
		handleRequest(r, req, requestBody, proxyURL)
	}

	// 执行代理
	proxy.ServeHTTP(w, r)
}

// 处理响应并存储到数据库
func handleResponse(r *http.Request, requestBody []byte, resp *http.Response) error {
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("error reading response body: %v", err)
		return err
	}
	resp.Body = io.NopCloser(bytes.NewBuffer(responseBody)) // 重置响应体供后续读取

	// 获取完整的请求路径（包括查询参数）
	fullRequestPath := r.URL.String()

	// 查询是否已存在相同路径和请求参数的记录
	var existingAPIResponse model.APIResponse
	err = component.DefaultDB().Where("endpoint = ? AND request_params = ?", fullRequestPath, string(requestBody)).First(&existingAPIResponse).Error
	if err == nil {
		// 更新现有记录
		existingAPIResponse.ResponseData = string(responseBody)
		existingAPIResponse.UpdatedAt = time.Now()

		if dbErr := component.DefaultDB().Save(&existingAPIResponse).Error; dbErr != nil {
			logrus.Errorf("failed to update existing response in database: %v", dbErr)
		} else {
			logrus.Infof("updated existing response in database: id=%d, path=%s", existingAPIResponse.ID, existingAPIResponse.Endpoint)
		}
	} else if err == gorm.ErrRecordNotFound {
		// 如果没有找到记录，创建新记录
		apiResponse := model.APIResponse{
			Endpoint:      fullRequestPath,
			RequestParams: string(requestBody),
			ResponseData:  string(responseBody),
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		if dbErr := component.DefaultDB().Create(&apiResponse).Error; dbErr != nil {
			logrus.Errorf("failed to store response in database: %v", dbErr)
		} else {
			logrus.Infof("stored response in database: id=%d, path=%s", apiResponse.ID, apiResponse.Endpoint)
		}
	} else {
		logrus.Errorf("failed to query database: %v", err)
	}

	logrus.Infof("response from target: status=%d, url=%s", resp.StatusCode, resp.Request.URL)
	return nil
}

// 处理代理错误，并从数据库返回历史数据
func handleProxyError(w http.ResponseWriter, req *http.Request, requestBody []byte, err error) {
	logrus.Errorf("proxy error: %v", err)

	// 如果代理错误，则查询数据库并返回历史数据
	fullRequestPath := req.URL.String()

	var existingAPIResponse model.APIResponse
	dbErr := component.DefaultDB().Where("endpoint = ? AND request_params = ?", fullRequestPath, string(requestBody)).First(&existingAPIResponse).Error
	if dbErr != nil {
		// 如果没有找到记录，返回代理错误
		http.Error(w, `{"error_msg":"proxy error"}`, http.StatusBadGateway)
		return
	}

	// 返回历史响应数据
	logrus.Infof("Returning cached response from database: id=%d, path=%s", existingAPIResponse.ID, existingAPIResponse.Endpoint)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(existingAPIResponse.ResponseData))
}

// 处理请求并设置代理目标
func handleRequest(r *http.Request, req *http.Request, requestBody []byte, proxyURL *url.URL) {
	// 设置目标地址
	req.URL.Scheme = proxyURL.Scheme
	req.URL.Host = proxyURL.Host
	req.Host = proxyURL.Host                              // 更新 Host Header
	req.Body = io.NopCloser(bytes.NewBuffer(requestBody)) // 复用请求体
	logrus.Infof("forwarding request: method=%s, path=%s, target=%s", req.Method, req.URL.Path, proxyURL)
}
