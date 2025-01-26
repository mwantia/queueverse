package ops

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/mwantia/prometheus/internal/agent/api"
	"github.com/mwantia/prometheus/internal/config"
	"github.com/mwantia/prometheus/internal/metrics"
	"github.com/mwantia/prometheus/internal/registry"
	"github.com/mwantia/prometheus/pkg/log"
)

type Server struct {
	Operation

	Log    log.Logger
	engine *gin.Engine
	srv    *http.Server
}

func (s *Server) Create(cfg *config.Config, reg *registry.PluginRegistry) (Cleanup, error) {
	s.Log = *log.New("server")
	s.engine = gin.Default()
	s.srv = &http.Server{
		Addr:    cfg.Server.Address,
		Handler: s.engine,
	}

	if err := s.addMiddlewares(); err != nil {
		return nil, fmt.Errorf("error adding routes: %w", err)
	}

	if err := s.addRoutes(cfg, reg); err != nil {
		return nil, fmt.Errorf("error adding routes: %w", err)
	}

	return func(ctx context.Context) error {
		return s.srv.Shutdown(ctx)
	}, nil
}

func (s *Server) Serve(ctx context.Context) error {
	s.Log.Info("Serving http server", "addr", s.srv.Addr)
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) addMiddlewares() error {
	s.engine.Use(func(c *gin.Context) {
		_, span := otel.Tracer("github.com/mwantia/prometheus").Start(c.Request.Context(), "Serve",
			trace.WithAttributes(
				attribute.String("method", c.Request.Method),
				attribute.String("addr", s.srv.Addr),
				attribute.String("fullpath", c.FullPath()),
			),
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()

		c.Next()
	})
	s.engine.Use(func(c *gin.Context) {
		observer := prometheus.ObserverFunc(func(v float64) {
			metrics.ServerHttpRequestsDurationSeconds.WithLabelValues(c.Request.Method, s.srv.Addr, c.FullPath()).Observe(v)
		})

		timer := prometheus.NewTimer(observer)
		defer timer.ObserveDuration()

		c.Next()

		metrics.ServerHttpRequestsTotal.WithLabelValues(c.Request.Method, s.srv.Addr, c.FullPath()).Inc()
	})

	return nil
}

func (s *Server) addRoutes(cfg *config.Config, reg *registry.PluginRegistry) error {
	v1 := s.engine.Group("/v1")
	auth := s.engine.Group("/v1", tokenAuthMiddleware(cfg.Server.Token))

	v1.GET("/health", api.HandleGetHealth(reg))
	v1.HEAD("/health", api.HandleIsHealthy(reg))

	auth.GET("/plugins", api.HandlePlugins(reg))
	auth.GET("/services", api.HandleServices(reg))

	auth.GET("/queue", api.HandleGetQueue(cfg))
	auth.GET("/queue/:task", api.HandleGetQueueTask(cfg))
	auth.HEAD("/queue/:task", api.HandleIsQueueTaskDone(cfg))
	auth.POST("/queue", api.HandlePostQueue(cfg))

	return nil
}
