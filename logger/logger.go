package logger

/**
  Customized logger module.
  Adding filename, line, request_id and user_id to prefix.
*/

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Logger struct {
	debugLogger   *log.Logger
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
}

type TraceDataKey struct{}

type TraceData struct {
	RequestId string
	UserId    uint64
}

func NewLogger() *Logger {
	logsFolderName := "logs"
	logFileName := "log"
	if !folderExists(logsFolderName) {
		err := os.MkdirAll(logsFolderName, os.ModePerm)
		if err != nil {
			return nil
		}
	}

	debugWriter, _ := os.OpenFile(filepath.Join(logsFolderName, logFileName), os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	infoWriter, _ := os.OpenFile(filepath.Join(logsFolderName, logFileName), os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	warningWriter, _ := os.OpenFile(filepath.Join(logsFolderName, logFileName), os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	errorWriter, _ := os.OpenFile(filepath.Join(logsFolderName, logFileName), os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)

	debugLogger := log.New(debugWriter, "[Debug] ", log.Llongfile|log.LstdFlags)
	infoLogger := log.New(infoWriter, "[Info] ", log.Llongfile|log.LstdFlags)
	warningLogger := log.New(warningWriter, "[Warning] ", log.Llongfile|log.LstdFlags)
	errorLogger := log.New(errorWriter, "[Error] ", log.Llongfile|log.LstdFlags)

	return &Logger{
		debugLogger:   debugLogger,
		infoLogger:    infoLogger,
		warningLogger: warningLogger,
		errorLogger:   errorLogger,
	}
}

// Debug for debug level
func (l *Logger) Debug(ctx context.Context, args ...interface{}) {
	_ = l.debugLogger.Output(2, l.getMsg(ctx, args))
}

// Info for info level
func (l *Logger) Info(ctx context.Context, args ...interface{}) {
	_ = l.infoLogger.Output(2, l.getMsg(ctx, args))
}

// Warning for warning level
func (l *Logger) Warning(ctx context.Context, args ...interface{}) {
	_ = l.warningLogger.Output(2, l.getMsg(ctx, args))
}

// Error for error level
func (l *Logger) Error(ctx context.Context, args ...interface{}) {
	_ = l.errorLogger.Output(2, l.getMsg(ctx, args))
}

func (l *Logger) getMsg(ctx context.Context, args ...interface{}) string {
	msg := fmt.Sprint(args...)
	msg = strings.Trim(msg, "[]")

	traceDataStr := ""
	traceData, ok := ctx.Value(TraceDataKey{}).(TraceData)
	if !ok {
		return msg
	}
	if traceData.RequestId != "" {
		traceDataStr = fmt.Sprintf("request_id: %v", traceData.RequestId)
	}
	if traceData.UserId != 0 {
		traceDataStr = fmt.Sprintf("%v user_id: %v", traceDataStr, traceData.UserId)
	}
	traceDataStr = strings.Trim(traceDataStr, " ")

	msg = fmt.Sprintf("{%v} %v", traceDataStr, msg)
	return msg
}

func folderExists(folderName string) bool {
	_, err := os.Stat(folderName)
	if err != nil && os.IsExist(err) {
		return true
	}
	return false
}
