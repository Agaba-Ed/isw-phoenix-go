package logger

import (
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

// InitLogger initializes the logger with file and console output
func InitLogger() {
	Log = logrus.New()

	// Set log level
	Log.SetLevel(logrus.InfoLevel)

	// Create logs directory if it doesn't exist
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		Log.WithError(err).Fatal("Failed to create logs directory")
	}

	// Create log file with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFile := filepath.Join(logsDir, "isw-phoenix-"+timestamp+".log")

	// Open log file
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		Log.WithError(err).Fatal("Failed to open log file")
	}

	// Set output to both file and console
	Log.SetOutput(file)

	// Also log to console
	Log.AddHook(&ConsoleHook{})

	// Set JSON formatter for structured logging
	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	Log.Info("Logger initialized successfully")
}

// ConsoleHook is a logrus hook to also output to console
type ConsoleHook struct{}

func (hook *ConsoleHook) Fire(entry *logrus.Entry) error {
	// Create a simple formatter for console output
	formatter := &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	}

	line, err := formatter.Format(entry)
	if err != nil {
		return err
	}

	os.Stdout.Write(line)
	return nil
}

func (hook *ConsoleHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// LogRequest logs HTTP request details
func LogRequest(method, url string, body interface{}, response interface{}, err error) {
	entry := Log.WithFields(logrus.Fields{
		"method":   method,
		"url":      url,
		"body":     body,
		"response": response,
	})

	if err != nil {
		entry.WithError(err).Error("Request failed")
	} else {
		entry.Info("Request completed")
	}
}

// LogAPIResponse logs API response with response code
func LogAPIResponse(responseCode string, response interface{}, endpoint string) {
	entry := Log.WithFields(logrus.Fields{
		"endpoint":     endpoint,
		"responseCode": responseCode,
		"response":     response,
	})

	// Log different levels based on response code
	switch responseCode {
	case "90000":
		entry.Info("API request successful")
	case "90063":
		entry.Warn("API request failed - Invalid request or authentication issue")
	default:
		entry.Warn("API request failed with unknown response code")
	}
}
