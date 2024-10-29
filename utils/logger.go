package utils

import (
	"fmt"
	"os"
	"time"
)

func LogOperation(logFile *os.File, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintf("[%s] %s\n", timestamp, message)
	fmt.Print(logMessage)
	logFile.WriteString(logMessage)
}
