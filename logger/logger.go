package logger

/**
  Self-defined logger module.
  Adding filename, line, request_id and user_id to prefix.
  Save files in folders named by current date.
*/

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type Logger struct {
	debugLogger   *log.Logger
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger

	date      string
	printKeys []string
}

func NewLogger() *Logger {
	date := time.Now().Format("2006-01-02")
	err := os.Mkdir(date, os.ModePerm)
	if err != nil {
		fmt.Println(err.Error())
	}

	debugWriter, _ := os.OpenFile("./logs/"+date+"/debug.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	infoWriter, _ := os.OpenFile("./logs/"+date+"/info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	warningWriter, _ := os.OpenFile("./logs/"+date+"/warning.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	errorWriter, _ := os.OpenFile("./logs/"+date+"/error.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)

	debugLogger := log.New(debugWriter, "[Debug] ", log.Lshortfile|log.LstdFlags)
	infoLogger := log.New(infoWriter, "[Info] ", log.Lshortfile|log.LstdFlags)
	warningLogger := log.New(warningWriter, "[Warning] ", log.Lshortfile|log.LstdFlags)
	errorLogger := log.New(errorWriter, "[Error] ", log.Lshortfile|log.LstdFlags)

	return &Logger{
		debugLogger:   debugLogger,
		infoLogger:    infoLogger,
		warningLogger: warningLogger,
		errorLogger:   errorLogger,
		date:          date,
		printKeys:     []string{"request_id", "user_id"},
	}
}

// Debug for debug level
func (l *Logger) Debug(ctx context.Context, args ...interface{}) {
	l.debugLogger.Println(l.getMsg(ctx, args))
}

// Info for info level
func (l *Logger) Info(ctx context.Context, args ...interface{}) {
	l.infoLogger.Println(l.getMsg(ctx, args))
}

// Warning for warning level
func (l *Logger) Warning(ctx context.Context, args ...interface{}) {
	l.warningLogger.Println(l.getMsg(ctx, args))
}

// Error for error level
func (l *Logger) Error(ctx context.Context, args ...interface{}) {
	l.errorLogger.Println(l.getMsg(ctx, args))
}

func (l *Logger) getMsg(ctx context.Context, args ...interface{}) string {
	msg := fmt.Sprint(args...)
	msg = strings.Trim(msg, "[]")
	header := ""
	for i, k := range l.printKeys {
		v := ctx.Value(k)
		if v == nil {
			continue
		}
		if i == 0 {
			header = fmt.Sprintf("%v: %v", k, v)
		} else {
			header = fmt.Sprintf("%v %v: %v", header, k, v)
		}
	}
	msg = fmt.Sprintf("{%v} %v", header, msg)
	return msg
}
