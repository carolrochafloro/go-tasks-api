package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var Logger *log.Logger

func NewLogger() {

	logDir := filepath.Join("app", "internal", "logging")
	logFilePath := filepath.Join(logDir, "app.log")

	// Cria o diretório se não existir
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		log.Fatal("Failed to create log directory:", err)
		panic("Failed to create log directory")
	}

	logFile, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatal("Failed to open log file:", err)
		panic("Failed to open log file")
	}
	Logger = log.New(logFile, "", log.LstdFlags|log.Lshortfile)

}

func Info(msg string, pairs ...interface{}) {
	Logger.Println("INFO", buildLogKeyValue(msg, pairs...))
}

func Warn(msg string, pairs ...interface{}) {
	Logger.Println("WARN", buildLogKeyValue(msg, pairs...))
}

func Error(msg string, pairs ...interface{}) {
	Logger.Println("ERROR", buildLogKeyValue(msg, pairs...))
}

func Debug(msg string, pairs ...interface{}) {
	Logger.Println("DEBUG", buildLogKeyValue(msg, pairs...))
}

func buildLogKeyValue(msg string, pairs ...interface{}) []interface{} {
	KeyValue := make([]interface{}, 0)
	KeyValue = append(KeyValue, parseBeforePrint("msg", msg))
	KeyValue = append(KeyValue, parseBeforePrint(pairs...))
	return KeyValue
}

func parseBeforePrint(r ...interface{}) (v string) {
	for i, e := range r {
		if i%2 == 0 {
			v += fmt.Sprintf(" %v=", e)
		} else {
			v += fmt.Sprintf("%v", e)
		}
	}
	return
}

