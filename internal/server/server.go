package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/nimishgj/aws-radar/internal/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

type Server struct {
	httpServer *http.Server
	config     *config.ServerConfig
}

func New(cfg *config.ServerConfig) *Server {
	mux := http.NewServeMux()

	// Metrics endpoint
	mux.Handle(cfg.MetricsPath, promhttp.Handler())

	// Health check endpoint
	mux.HandleFunc(cfg.HealthPath, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Root endpoint
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head><title>AWS Radar</title></head>
<body>
<h1>AWS Radar</h1>
<p>AWS Resource Monitoring Agent</p>
<ul>
<li><a href="%s">Metrics</a></li>
<li><a href="%s">Health</a></li>
</ul>
</body>
</html>
`, cfg.MetricsPath, cfg.HealthPath)))
	})

	return &Server{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Port),
			Handler:      mux,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		config: cfg,
	}
}

func (s *Server) Start() error {
	log.Info().
		Int("port", s.config.Port).
		Str("metrics_path", s.config.MetricsPath).
		Str("health_path", s.config.HealthPath).
		Msg("Starting HTTP server")

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Info().Msg("Shutting down HTTP server")
	return s.httpServer.Shutdown(ctx)
}
