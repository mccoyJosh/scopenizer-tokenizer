package util

import (
	"fmt"
	"time"
)

func PrintLog(logType string, message string) {
	fmt.Printf("%s: ", logType)
	PrintTime()
	fmt.Printf("\tNOTE: %s\n", message)
}

func PrintTime() {
	t := time.Now()
	fmt.Printf("%d/%d/%d  ", t.Year(), t.Month(), t.Day())
	fmt.Printf("%d:%d:%d.%d ", t.Hour(), t.Minute(), t.Second(), t.Nanosecond())
}

func Fatal(msg string) {
	PrintLog("FATAL", msg)
	panic(msg)
}

func Error(msg string, err error) {
	PrintLog("ERROR", msg)
	if err != nil {
		fmt.Print("ERROR MESSAGE: ")
		fmt.Println(err.Error())
	}
}

func Warning(msg string) {
	PrintLog("WARNING", msg)
}

func Info(msg string) {
	PrintLog("INFO", msg)
}

func Debug(msg string) {
	PrintLog("DEBUG", msg)
}
