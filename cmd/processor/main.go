package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/tartoide/stori/stori-challenge/internal/config"
	"github.com/tartoide/stori/stori-challenge/internal/infrastructure/database"
	"github.com/tartoide/stori/stori-challenge/internal/infrastructure/email"
	"github.com/tartoide/stori/stori-challenge/internal/infrastructure/file"
	"github.com/tartoide/stori/stori-challenge/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Setup logger
	logger := setupLogger(cfg.Server.LogLevel)
	logger.Info("starting Stori transaction processor",
		"environment", cfg.Server.Environment,
		"log_level", cfg.Server.LogLevel)

	// Create context that listens for interrupt signals
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Initialize services
	app, cleanup, err := initializeApplication(cfg, logger)
	if err != nil {
		logger.Error("failed to initialize application", "error", err)
		os.Exit(1)
	}
	defer cleanup()

	// Start processing
	recipientEmail := getRecipientEmail(cfg)
	logger.Info("starting file processing",
		"watch_directory", cfg.File.WatchDirectory,
		"recipient", recipientEmail)

	if err := app.WatchAndProcess(ctx, cfg.File.WatchDirectory, recipientEmail); err != nil {
		if err == context.Canceled {
			logger.Info("application stopped gracefully")
		} else {
			logger.Error("application error", "error", err)
			os.Exit(1)
		}
	}
}

func setupLogger(level string) *slog.Logger {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	return slog.New(handler)
}

func initializeApplication(cfg *config.Config, logger *slog.Logger) (*services.TransactionProcessor, func(), error) {
	var cleanupFuncs []func()

	cleanup := func() {
		for i := len(cleanupFuncs) - 1; i >= 0; i-- {
			cleanupFuncs[i]()
		}
	}

	// Initialize file processor
	fileProcessor := file.NewCSVFileProcessor(logger)

	// Initialize email service
	emailService := email.NewSMTPEmailService(cfg.Email, logger)

	// Initialize summary calculator
	calculator := services.NewSummaryCalculator()

	// Initialize database (optional)
	var dataStore services.DataStore
	if cfg.DatabaseEnabled() {
		logger.Info("initializing database connection")
		store, err := database.NewPostgresDataStore(cfg.Database, logger)
		if err != nil {
			logger.Warn("database initialization failed, continuing without database", "error", err)
		} else {
			dataStore = store
			if closer, ok := store.(interface{ Close() error }); ok {
				cleanupFuncs = append(cleanupFuncs, func() {
					if err := closer.Close(); err != nil {
						logger.Error("failed to close database connection", "error", err)
					}
				})
			}
			logger.Info("database connection established successfully")
		}
	} else {
		logger.Info("database disabled, continuing without persistence")
	}

	// Create transaction processor
	processor := services.NewTransactionProcessor(
		fileProcessor,
		emailService,
		calculator,
		dataStore,
		logger,
	)

	return processor, cleanup, nil
}

func getRecipientEmail(cfg *config.Config) string {
	// Check for environment variable first
	if email := os.Getenv("RECIPIENT_EMAIL"); email != "" {
		return email
	}

	// Fall back to a default for demo purposes
	if cfg.IsDevelopment() {
		return "demo@stori.com"
	}

	// In production, this should be required
	return "notifications@stori.com"
}
