package configure

import (
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

var logger *slog.Logger

func SetupLogging() *slog.Logger {
	// Create logs directory if it doesn't exist
	logDir := "logs"
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		slog.Error("Failed to create log directory", "error", err)
		return nil
	}

	// Create log file
	logFilePath := filepath.Join(logDir, time.Now().Format("2006-01-02")+".log")
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		slog.Error("Failed to open log file", "error", err)
		return nil
	}

	// Create a new logger with a JSON handler
	handlerOptions := &slog.HandlerOptions{}
	handler := slog.NewJSONHandler(logFile, handlerOptions)
	logger = slog.New(handler)

	// Log the setup completion
	logger.Info("Logging setup complete")
	return logger
}
