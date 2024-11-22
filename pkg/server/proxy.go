package server

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"rbd_proxy_dp/config"
	"rbd_proxy_dp/pkg/component"
	"rbd_proxy_dp/pkg/gogo"
	"syscall"
	"time"
)

type Server struct {
	ctx        context.Context
	cancel     context.CancelFunc
	components []component.Component
	serverName string
}

func New(ctx context.Context) *Server {
	ctx, cancel := context.WithCancel(ctx)
	cago := &Server{
		ctx:        ctx,
		cancel:     cancel,
		serverName: config.DefaultPublic().ServerName,
	}
	return cago
}

func (r *Server) Registry(component component.Component) *Server {
	err := component.Start(r.ctx)
	if err != nil {
		panic(err)
	}
	r.components = append(r.components, component)
	return r
}

func (r *Server) Start() error {
	// 创建信号监听通道
	quitSignal := make(chan os.Signal, 1)
	signal.Notify(
		quitSignal,
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer signal.Stop(quitSignal) // 确保退出时释放信号通知

	// 等待信号或上下文取消
	select {
	case <-quitSignal:
		logrus.Infof("received termination signal")
		r.cancel()
	case <-r.ctx.Done():
		logrus.Infof("context cancelled")
	}

	// 关闭所有组件
	logrus.Infof("begin close all component")
	for _, cpt := range r.components {
		cpt.CloseHandle()
	}

	// 等待所有异步任务完成
	stopCh := make(chan struct{})
	go func() {
		gogo.Wait() // 假设 gogo.Wait() 是一个阻塞方法，等待所有任务完成
		close(stopCh)
	}()

	select {
	case <-stopCh:
		logrus.Infof("all components stopped gracefully")
	case <-time.After(10 * time.Second):
		logrus.Warnf("timeout waiting for components to stop")
	}

	logrus.Infof("%v server has stopped", r.serverName)
	return nil
}
