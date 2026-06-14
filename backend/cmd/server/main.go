package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"oa-nsdiy/backend/internal/config"
	"oa-nsdiy/backend/internal/db"
	"oa-nsdiy/backend/internal/pkg/logger"
	"oa-nsdiy/backend/internal/repository"
	"oa-nsdiy/backend/internal/service"
)

func seedAdminUser() {
	ctx := context.Background()
	userRepo := repository.NewUserRepository(db.Client)
	authService := service.NewAuthService(userRepo, "", 0, 0)

	// Check if admin already exists
	_, err := userRepo.GetByUsername(ctx, "admin")
	if err == nil {
		return
	}

	salt, _ := authService.GenerateSalt(ctx)
	hashedPassword, _ := authService.HashPassword(ctx, "admin123", salt)

	nickname := "Administrator"
	admin := &repository.User{
		Username:       "admin",
		Email:          "admin@example.com",
		Nickname:       &nickname,
		Salt:           salt,
		HashedPassword: hashedPassword,
		UserType:       "HUMAN",
		IsActive:       true,
		TokenVersion:   1,
	}

	if err := userRepo.Create(ctx, admin); err != nil {
		logger.S().Warnw("Failed to seed admin user", "error", err)
		return
	}

	logger.S().Infow("Default admin user created", "username", "admin", "password", "admin123")
}

func main() {
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

	seedAdminUser()

	// Initialize application using Wire dependency injection
	app, err := initializeApplication(db.Client, cfg)
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
