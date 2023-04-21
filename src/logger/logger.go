package logger

import (
	"fmt"
	"time"
)

const (
	reset  = "\x1B[0m"
	red    = "\x1B[31m"
	green  = "\x1B[32m"
	yellow = "\x1B[93m"
	bold   = "\x1B[01m"
)

type Logger interface {
	Info(f string, args ...interface{})
	Warning(f string, args ...interface{})
	Error(f string, args ...interface{})
	Critical(f string, args ...interface{})
}

type ColorLogger struct {
	Name string
}

func NewColorLogger(name string) ColorLogger {
	return ColorLogger{Name: name}
}

func (l ColorLogger) getTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func (l ColorLogger) Info(f string, args ...interface{}) {
	fmt.Printf(fmt.Sprintf("[%s] [%s/%s]: %s\n", l.getTime(), l.Name, l.colorText("INFO", "green", false), f), args...)
}

func (l ColorLogger) Warning(f string, args ...interface{}) {
	fmt.Printf(fmt.Sprintf("[%s] [%s/%s]: %s\n", l.getTime(), l.Name, l.colorText("WARNING", "yellow", true), l.colorText(f, "yellow", false)), args...)
}

func (l ColorLogger) Error(f string, args ...interface{}) {
	fmt.Printf(fmt.Sprintf("[%s] [%s/%s]: %s\n", l.getTime(), l.Name, l.colorText("ERROR", "red", false), l.colorText(f, "red", false)), args...)
}

func (l ColorLogger) Critical(f string, args ...interface{}) {
	fmt.Printf(fmt.Sprintf("[%s] [%s/%s]: %s\n", l.getTime(), l.Name, l.colorText("CRITICAL", "red", true), l.colorText(f, "red", true)), args...)
}

func (l ColorLogger) colorText(text, color string, isBold bool) string {
	var formatted string
	switch color {
	case "green":
		formatted += green
	case "red":
		formatted += red
	case "yellow":
		formatted += yellow
	}
	if isBold {
		formatted += bold
	}
	formatted = formatted + text
	return formatted + reset
}
