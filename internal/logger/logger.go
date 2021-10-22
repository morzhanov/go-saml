package logger

import (
	"fmt"
	"os"

	lgr "github.com/crewjam/saml/logger"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() (*zap.Logger, error) {
	baseLogsPath := "./logs/"
	if err := os.MkdirAll(baseLogsPath, 0777); err != nil {
		return nil, err
	}

	// info
	infoFilePath := fmt.Sprintf("%sinfo.log", baseLogsPath)
	_, err := os.Create(infoFilePath)
	if err != nil {
		return nil, err
	}
	infoF, err := os.OpenFile(infoFilePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return nil, err
	}

	// error
	errFilePath := fmt.Sprintf("%serror.log", baseLogsPath)
	_, err = os.Create(errFilePath)
	if err != nil {
		return nil, err
	}
	errF, err := os.OpenFile(errFilePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return nil, err
	}

	fileEncoder := zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(infoF), zap.InfoLevel),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(errF), zap.ErrorLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.InfoLevel),
	)
	return zap.New(core), nil
}

type idpLogger struct {
	log *zap.Logger
}

func (l *idpLogger) Printf(format string, v ...interface{}) { l.log.Printf(format, v) }
func (l *idpLogger) Print(v ...interface{})                 { l.log.Print(v) }
func (l *idpLogger) Println(v ...interface{})               { l.log.Println(v) }
func (l *idpLogger) Fatal(v ...interface{})                 { l.log.Fatal(v) }
func (l *idpLogger) Fatalf(format string, v ...interface{}) { l.log.Fatalf(format, v) }
func (l *idpLogger) Fatalln(v ...interface{})               { l.log.Fatalln(v) }
func (l *idpLogger) Panic(v ...interface{})                 { l.log.Panic(v) }
func (l *idpLogger) Panicf(format string, v ...interface{}) { l.log.Panicf(format, v) }
func (l *idpLogger) Panicln(v ...interface{})               { l.log.Panicln(v) }

func NewIDPLogger(log *zap.Logger) lgr.Interface {
	return &idpLogger{log}
}
