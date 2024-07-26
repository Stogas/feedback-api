package main

import (
	"fmt"
	"log/slog"

	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/Stogas/feedback-api/internal/config"
)

func initMetrics() (*gin.Engine, *ginprom.Prometheus) {
	r := gin.Default()

	p := ginprom.New(
		ginprom.Engine(r),
		ginprom.Subsystem("feedbackapi"),
		ginprom.Path("/metrics"),
	)

	return r, p
}

func startMetrics(r *gin.Engine, conf config.MetricsConfig) {
	slog.Info("Starting Prometheus exporter", "host", conf.Host, "port", conf.Port)

	if err := r.Run(fmt.Sprintf("%s:%v", conf.Host, conf.Port)); err != nil {
		slog.Error("Failed to run metrics exporter", "error", err)
	}
}
