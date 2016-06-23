package main

import (
	"io"
	"log"
	"os"
)

type ServerLogger struct {
	Warning *log.logger
}

func NewServerLogger() *ServerLogger {
	logfile, err := os.OpenFile("log-output.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file:", err)
	}

	return &ServerLogger{
		Warning: log.New(io.Multiwriter(logfile, os.Stdout), "[WARNING] ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
