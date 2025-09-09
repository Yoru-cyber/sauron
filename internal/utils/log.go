package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var (
	ProgramStartTime time.Time
	logFile          *os.File
	isInitialized    bool
)

func InitLogger() error {
	if isInitialized {
		return nil
	}

	ProgramStartTime = time.Now()

	// Create logs directory
	if err := os.MkdirAll("logs", 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Create log file with timestamp
	filename := fmt.Sprintf("logs/app_%s.log",
		ProgramStartTime.Format("2006-01-02_15-04-05"))

	var err error
	logFile, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	// Set log output to both file and stdout
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	isInitialized = true
	log.Printf("Logger initialized at: %v", ProgramStartTime)

	return nil
}

func CleanupLogger() {
	if logFile != nil {
		duration := time.Since(ProgramStartTime)
		log.Printf("Program finished. Duration: %v", duration)
		logFile.Close()
	}
}

// Log functions that work across packages
func LogInfo(format string, v ...interface{}) {
	log.Printf("[INFO] "+format, v...)
}

func LogError(format string, v ...interface{}) {
	log.Printf("[ERROR] "+format, v...)
}

func LogDebug(format string, v ...interface{}) {
	log.Printf("[DEBUG] "+format, v...)
}
