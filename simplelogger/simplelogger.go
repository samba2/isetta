package simplelogger

// after a lot of investigation decided to roll my own simplelogger.
// wanted basic logger + log levels - nothing more, so there you go.
import (
	"fmt"
	"log"
	"os"

	"golang.org/x/exp/maps"
)

// LogLevel type
type LogLevel int

const (
	// LevelTrace logs everything
	LevelTrace LogLevel = (1 << iota)

	// LevelDebug logs Debug and above
	LevelDebug

	// LevelInfo logs Info and above
	LevelInfo

	// LevelWarn logs Warning and Errors
	LevelWarn

	// LevelError logs just Errors
	LevelError
)

var Levels = map[string]LogLevel{
	"trace": LevelTrace,
	"debug": LevelDebug,
	"info": LevelInfo,
	"warn": LevelWarn,
	"error": LevelError,
}

// global logger
var Logger = NewSimpleLogger(LevelTrace)

type Sl struct {
	CurrentLogLevel LogLevel
	logger *log.Logger
}

func GetValidLogLevels() []string {
	return maps.Keys(Levels)
}

func NewSimpleLogger(logLevel LogLevel) Sl {
	return Sl{
		CurrentLogLevel: logLevel,
		logger: log.New(os.Stdout, "", 0),
	}
}

func (sl Sl) Trace(format string, v ...any) {
	if sl.CurrentLogLevel <= LevelTrace {
		sl.logger.Println("Trace:", buildOutput(format, v...))
	}
}

func (sl Sl)Debug(format string, v ...any) {
	if sl.CurrentLogLevel <= LevelDebug {
		sl.logger.Println("Debug:", buildOutput(format, v...))
	}
}

func (sl Sl)Info(format string, v ...any) {
	if sl.CurrentLogLevel <= LevelInfo {
		sl.logger.Println("Info:", buildOutput(format, v...))
	}
}

func (sl Sl)Warn(format string, v ...any) {
	if sl.CurrentLogLevel <= LevelWarn {
		sl.logger.Println("Warn:", buildOutput(format, v...))
	}
}

func (sl Sl)Error(format string, v ...any) {
	if sl.CurrentLogLevel <= LevelError {
		sl.logger.Fatalln("Fatal:", buildOutput(format, v...))
	}
}

func (sl Sl)Error2(err error) {
	if sl.CurrentLogLevel <= LevelError {
		sl.logger.Fatalln("Fatal:", err)
	}
}


func buildOutput(format string, v ...any)  string {
	return fmt.Sprintf(format, v...)
}