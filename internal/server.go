package internal

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/dchlong/billing-be/internal/config"
	dhttp "github.com/dchlong/billing-be/internal/delivery/http"
	"github.com/dchlong/billing-be/internal/middlewares"
	"github.com/dchlong/billing-be/pkg/logger"
)

type Server struct {
	appCfg  *config.AppConfig
	handler dhttp.Handler
	router  *gin.Engine
	logger  logger.ILogger
}

func ProvideRouter(ilogger logger.ILogger) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middlewares.SetRequestID())
	router.Use(middlewares.SetupLog(ilogger))
	return router
}

func ProvideServer(
	appCfg *config.AppConfig, router *gin.Engine, handler dhttp.Handler, ilogger logger.ILogger,
) (*Server, func(), error) {
	handler.Register(router)
	return &Server{
			appCfg:  appCfg,
			router:  router,
			handler: handler,
			logger:  ilogger,
		}, func() {
		}, nil
}

func (s *Server) Start() error {
	httpServer := &http.Server{
		Handler:           s.router,
		Addr:              s.appCfg.HTTPAddr,
		ReadHeaderTimeout: time.Duration(s.appCfg.ReadHeaderTimeout) * time.Millisecond,
	}

	s.logger.Infof("starting server at port %s", s.appCfg.HTTPAddr)
	return httpServer.ListenAndServe()
}
