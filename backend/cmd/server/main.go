package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"oa-nsdiy/backend/internal/config"
	"oa-nsdiy/backend/internal/db"
	"oa-nsdiy/backend/internal/pkg/logger"
)

//go:embed VERSION
var embeddedVersion string

var (
	Version   = "dev"
	Commit    = ""
	Date      = ""
	BuildType = ""
)

func init() {
	if v := strings.TrimSpace(embeddedVersion); v != "" {
		Version = v
	}
}

func main() {
	showVersion := flag.Bool("version", false, "Show version information")
	flag.Parse()

	if *showVersion {
		log.Printf("oa-nsdiy %s (commit: %s, built: %s, type: %s)\n", Version, Commit, Date, BuildType)
		return
	}

	cfg, err := config.Load()
	if err != nil {
		logger.S().Fatalw("Failed to load config", "error", err)
	}

	if err := logger.Init(logger.Options{
		Level:       cfg.Log.Level,
		Format:      cfg.Log.Format,
		Caller:      cfg.Log.Caller,
		ServiceName: cfg.Log.ServiceName,
		Output: logger.OutputOptions{
			ToStdout: cfg.Log.Output.ToStdout,
			ToFile:   cfg.Log.Output.ToFile,
			FilePath: cfg.Log.Output.FilePath,
		},
		Rotation: logger.RotationOptions{
			MaxSizeMB:  cfg.Log.Rotation.MaxSizeMB,
			MaxBackups: cfg.Log.Rotation.MaxBackups,
			MaxAgeDays: cfg.Log.Rotation.MaxAgeDays,
			Compress:   cfg.Log.Rotation.Compress,
		},
	}); err != nil {
		logger.S().Fatalw("Failed to initialize logger", "error", err)
	}
	defer logger.Sync()

	if err := db.Init(cfg.Database.Driver, cfg.Database.Source); err != nil {
		logger.S().Fatalw("Failed to initialize database", "error", err)
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		logger.S().Fatalw("Failed to run migrations", "error", err)
	}

	// Initialize application using Wire dependency injection
	app, err := initializeApplication(db.Client, cfg, Version)
	if err != nil {
		logger.S().Fatalw("Failed to initialize application", "error", err)
	}
	defer app.Cleanup()

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.S().Infow("Server starting", "address", addr)

	httpSrv := &http.Server{Addr: addr, Handler: app.Router}

	go func() {
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.S().Fatalw("Failed to start server", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.S().Infow("Server shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpSrv.Shutdown(ctx); err != nil {
		logger.S().Errorw("Server forced to shutdown", "error", err)
	}

	logger.S().Infow("Server exited")
}
