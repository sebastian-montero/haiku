package logger

import "log"

func Info(message string) {
	log.Println(message)
}

func Error(message string) {
	log.Println("[ERROR]", message)
}
