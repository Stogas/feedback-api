package main

import (
	"fmt"
	"log/slog"

	"github.com/Depado/ginprom"
	"github.com/Stogas/feedback-api/internal/config"
	"github.com/gin-gonic/gin"
)

func initMetrics(m []gin.HandlerFunc) (*gin.Engine, *ginprom.Prometheus) {
	r := gin.New()

	r.Use(gin.Recovery())
	for _, m := range m {
		// slog.Debug("Gin: Adding middleware")
		r.Use(m)
	}

	p := ginprom.New(
		ginprom.Engine(r),
		ginprom.Subsystem("feedbackapi"),
		ginprom.Path("/metrics"),
	)

	p.AddCustomCounter("satisfaction", "Counts how many good/bad satisfactions are received successfully. Note that this only counts new submittions, not updates", []string{"satisfied"})

	return r, p
}

func startMetrics(r *gin.Engine, conf config.MetricsConfig) {
	slog.Info("Starting Prometheus exporter", "host", conf.Host, "port", conf.Port)

	if err := r.Run(fmt.Sprintf("%s:%v", conf.Host, conf.Port)); err != nil {
		slog.Error("Failed to run metrics exporter", "error", err)
	}
}
