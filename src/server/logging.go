// logging @exprays
// Copyright (c) 2023, exprays <
// In future more logging levels will be added
// License: MIT

package server

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	// Logger instances
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
	CommandLogger *log.Logger

	// File handles to close
	logFiles []*os.File

	// Mutex for thread-safe logging
	logMutex sync.Mutex
)

const (
	LogDir         = "logs"
	InfoLogFile    = "info.log"
	ErrorLogFile   = "error.log"
	CommandLogFile = "commands.log"
)

// InitLogging sets up the logging system
func InitLogging() error {
	logMutex.Lock()
	defer logMutex.Unlock()

	// Create logs directory if it doesn't exist
	if err := os.MkdirAll(LogDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	// Set up info logger
	infoFile, err := os.OpenFile(filepath.Join(LogDir, InfoLogFile), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open info log file: %v", err)
	}
	logFiles = append(logFiles, infoFile)
	InfoLogger = log.New(infoFile, "INFO: ", log.Ldate|log.Ltime)

	// Set up error logger
	errorFile, err := os.OpenFile(filepath.Join(LogDir, ErrorLogFile), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open error log file: %v", err)
	}
	logFiles = append(logFiles, errorFile)
	ErrorLogger = log.New(errorFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	// Set up command logger
	cmdFile, err := os.OpenFile(filepath.Join(LogDir, CommandLogFile), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open command log file: %v", err)
	}
	logFiles = append(logFiles, cmdFile)
	CommandLogger = log.New(cmdFile, "", log.Ldate|log.Ltime)

	// Log startup
	InfoLogger.Println("Logging system initialized")

	// Set up log rotation
	go rotateLogsDaily()

	return nil
}

// CloseLogFiles closes all log files
func CloseLogFiles() {
	logMutex.Lock()
	defer logMutex.Unlock()

	for _, file := range logFiles {
		file.Close()
	}
	logFiles = nil
}

// LogInfo logs an info message
func LogInfo(format string, v ...interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()

	if InfoLogger != nil {
		InfoLogger.Printf(format, v...)
		fmt.Printf("[INFO] "+format+"\n", v...)
	}
}

// LogError logs an error message
func LogError(format string, v ...interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()

	if ErrorLogger != nil {
		ErrorLogger.Printf(format, v...)
		fmt.Printf("[ERROR] "+format+"\n", v...)
	}
}

// LogCommand logs a command
func LogCommand(clientIP string, command string) {
	logMutex.Lock()
	defer logMutex.Unlock()

	if CommandLogger != nil {
		CommandLogger.Printf("[%s] %s", clientIP, command)
	}
}

// rotateLogsDaily rotates logs once per day
func rotateLogsDaily() {
	for {
		// Calculate time until next day
		now := time.Now()
		nextDay := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		timeUntilNextDay := nextDay.Sub(now)

		// Sleep until next day
		time.Sleep(timeUntilNextDay)

		// Rotate logs
		rotateLogFiles()
	}
}

// rotateLogFiles rotates all log files
func rotateLogFiles() {
	logMutex.Lock()
	defer logMutex.Unlock()

	timestamp := time.Now().Format("2006-01-02")

	// Close existing files
	for _, file := range logFiles {
		file.Close()
	}
	logFiles = nil

	// Rename existing log files with timestamp
	renameLogFile(InfoLogFile, timestamp)
	renameLogFile(ErrorLogFile, timestamp)
	renameLogFile(CommandLogFile, timestamp)

	// Reinitialize logging
	_ = InitLogging() // Ignore errors during rotation
}

// renameLogFile renames a log file with a timestamp
func renameLogFile(fileName, timestamp string) {
	oldPath := filepath.Join(LogDir, fileName)
	newPath := filepath.Join(LogDir, fmt.Sprintf("%s.%s", fileName, timestamp))

	// If the file exists, rename it
	if _, err := os.Stat(oldPath); err == nil {
		if err := os.Rename(oldPath, newPath); err != nil {
			fmt.Printf("Error renaming log file %s: %v\n", fileName, err)
		}
	}
}
