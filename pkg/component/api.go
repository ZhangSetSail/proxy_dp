package component

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"rbd_proxy_dp/config"
	"rbd_proxy_dp/pkg/gogo"
	"time"
)

type API struct {
	mux    *http.ServeMux
	server *http.Server
}

func NewAPI(register func(mux *http.ServeMux)) *API {
	mux := http.NewServeMux()
	register(mux)
	return &API{
		mux: mux,
	}
}

func (r *API) Start(ctx context.Context) error {
	return gogo.Go(ctx, func() error {
		for {
			select {
			case <-ctx.Done():
				return nil
			default:
				logrus.Infof("listen :%v", config.DefaultPublic().Port)
				logrus.Infof("%v", config.DefaultPublic().Port)

				r.server = &http.Server{
					Addr:    fmt.Sprintf(":%v", config.DefaultPublic().Port),
					Handler: r.mux,
				}
				if err := r.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logrus.Errorf("HTTP server ListenAndServe: %v", err)
				}
			}
		}
	})
}

func (r *API) CloseHandle() {
	if r.server != nil {
		logrus.Infof("begin close component http server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := r.server.Shutdown(shutdownCtx); err != nil {
			logrus.Errorf("http server shutdown: %v", err)
		} else {
			logrus.Infof("http server stopped gracefully.")
		}
	}
}
