package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"gorm.io/gorm/logger"
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
		ServiceName: "product-catalog-service",
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
		return
	}

	conn, err := net.Dial("tcp", "logstash-service:5044")
	if err != nil {
		Error("Connecting to Logstash error: %v\n", err)
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
	serviceDir := "product-catalog-service"
	logDir := filepath.Join(appDir, serviceDir)

	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("Log dizini oluşturulamadı: %v", err)
		return err
	}

	logFileName := fmt.Sprintf("product-catalog-service-%s.log", time.Now().Format("2006-01-02"))
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

func Info(format string, v ...interface{}) {
	log.Printf("INFO: "+format, v...)
}

func Error(format string, v ...interface{}) {
	log.Printf("ERROR: "+format, v...)
}

type GormLogger struct{}

func (g GormLogger) LogMode(logger.LogLevel) logger.Interface {
	// Bu örnekte log seviyesi değişikliklerini yönetmiyoruz ama gerekiyorsa burada yönetebilirsiniz.
	return g
}

func (GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	Log("INFO", fmt.Sprintf(msg, data...))
}

func (GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	Log("WARN", fmt.Sprintf(msg, data...))
}

func (GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	Log("ERROR", fmt.Sprintf(msg, data...))
}

func (GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if err != nil {
		sql, rows := fc()
		Log("ERROR", fmt.Sprintf("%s [%d rows affected or returned] [%v]", sql, rows, err))
	} else {
		sql, rows := fc()
		Log("INFO", fmt.Sprintf("%s [%d rows affected or returned]", sql, rows))
	}
}
