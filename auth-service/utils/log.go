package utils

import (
	"auth-service/config"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"
)

var LogFile *os.File

type LogMessage struct {
	ServiceName string `json:"service_name"`
	Level       string `json:"level"`
	Message     string `json:"message"`
}

func Log(level, message string) {
	logMsg := LogMessage{
		ServiceName: config.Configuration.ServiceName,
		Level:       level,
		Message:     message,
	}
	logMsgJson := ToJSON(logMsg)
	if logMsg.Level == "INFO" {
		Info(logMsgJson)
	} else {
		Error(logMsgJson)
	}
	sendToLogstash(logMsg)
}
func sendToLogstash(logMsg LogMessage) {
	msgBytes, err := json.Marshal(logMsg)
	if err != nil {
		Error("Log message JSON marshal error: %v\n", err)
		//fmt.Printf("Log message JSON marshal error: %v\n", err)
		return
	}

	conn, err := net.Dial("tcp", "logstash-service:5044")
	if err != nil {
		Error("Connecting to Logstash error: %v\n", err)
		//fmt.Printf("Connecting to Logstash error: %v\n", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write(msgBytes)
	if err != nil {
		Error("Sending log to Logstash error: %v\n", err)
		//fmt.Printf("Sending log to Logstash error: %v\n", err)
	}
}
func LogFileInit() error {
	appDir := "../logs"
	serviceDir := "auth-service"
	logDir := filepath.Join(appDir, serviceDir)

	if err := os.MkdirAll(logDir, 0755); err != nil {

		log.Fatalf("Log dizini oluşturulamadı: %v", err)
		return err
	}

	logFileName := fmt.Sprintf("auth-service-%s.log", time.Now().Format("2006-01-02"))
	logFilePath := filepath.Join(logDir, logFileName)
	LogFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Log dosyası açılamadı: %v", err)
		return err
	}
	log.SetOutput(LogFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	return nil
}

// Info log seviyesi için bir wrapper fonksiyon
func Info(format string, v ...interface{}) {
	log.Printf("INFO: "+format, v...)
}

// Error log seviyesi için bir wrapper fonksiyon
func Error(format string, v ...interface{}) {
	log.Printf("ERROR: "+format, v...)
}

//
//package utils
//
//import (
//	"encoding/json"
//	"fmt"
//	"log"
//	"net"
//	"os"
//	"path/filepath"
//	"time"
//)
//
//type LogMessage struct {
//	ServiceName string `json:"service_name"`
//	Level       string `json:"level"`
//	Message     string `json:"message"`
//}
//
//// Logging encapsulates logging functionality with support for writing to files and sending logs to Logstash.
//type Logging struct {
//	LogFile     *os.File
//	ServiceName string
//}
//
//// NewLogger initializes a new logger instance, creating a log file in the process.
//func NewLogger(serviceName string) (*Logging, error) {
//	logFilePath, err := prepareLogFile()
//	if err != nil {
//		return nil, err
//	}
//
//	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
//	if err != nil {
//		return nil, fmt.Errorf("failed to open log file: %w", err)
//	}
//
//	log.SetOutput(file)
//	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
//
//	return &Logging{
//		LogFile:     file,
//		ServiceName: serviceName,
//	}, nil
//}
//
//// prepareLogFile creates the directory structure for the log files and returns the path to the log file.
//func prepareLogFile() (string, error) {
//	logDir := filepath.Join("../logs", "auth-service")
//	if err := os.MkdirAll(logDir, 0755); err != nil {
//		return "", fmt.Errorf("failed to create log directory: %w", err)
//	}
//
//	logFileName := fmt.Sprintf("auth-service-%s.log", time.Now().Format("2006-01-02"))
//	return filepath.Join(logDir, logFileName), nil
//}
//
//// Log writes a log message with the given level.
//func (l *Logging) Log(level, message string) {
//	logMsg := LogMessage{
//		ServiceName: l.ServiceName,
//		Level:       level,
//		Message:     message,
//	}
//
//	logMsgJson, _ := json.Marshal(logMsg) // Error handling omitted for brevity.
//	log.Printf("%s: %s\n", level, logMsgJson)
//
//	if err := l.sendToLogstash(logMsg); err != nil {
//		log.Printf("ERROR: failed to send log to Logstash: %v", err)
//	}
//}
//
//// sendToLogstash sends a log message to Logstash.
//func (l *Logging) sendToLogstash(logMsg LogMessage) error {
//	conn, err := net.Dial("tcp", "localhost:5044")
//	if err != nil {
//		return fmt.Errorf("connecting to Logstash error: %w", err)
//	}
//	defer conn.Close()
//
//	msgBytes, _ := json.Marshal(logMsg) // Error handling omitted for brevity.
//	if _, err := conn.Write(msgBytes); err != nil {
//		return fmt.Errorf("sending log to Logstash error: %w", err)
//	}
//
//	return nil
//}
