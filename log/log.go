package log

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Log struct {
	logger *zap.SugaredLogger
}

func New(logname, loglevel string) *Log {
	return &Log{logger: initLogger(filepath.Join("log", logname+".log"), loglevel)}
}

func initLogger(logPath string, loglevel string) *zap.SugaredLogger {
	hook := lumberjack.Logger{
		Filename:   logPath, // 日志文件路径，默认 os.TempDir()
		MaxSize:    10,      // 每个日志文件保存10M，默认 100M
		MaxBackups: 30,      // 保留30个备份，默认不限
		MaxAge:     7,       // 保留7天，默认不限
		Compress:   true,    // 是否压缩，默认不压缩
	}
	write := zapcore.AddSync(&hook)
	var level zapcore.Level
	switch loglevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "lineNum",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel: func(level zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(fmt.Sprintf("[%s]", level.CapitalString()))
		},
		EncodeTime:     zapcore.TimeEncoderOfLayout("[jt] 2006-01-02 15:04:05.000000"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(write)),
		level,
	)
	caller := zap.AddCaller()
	development := zap.Development()
	logger := zap.New(core, caller, development, zap.AddCallerSkip(1))
	return logger.Sugar()
}

func (l *Log) Debug(args ...interface{}) {
	l.logger.Debug(args)
}

func (l *Log) Info(args ...interface{}) {
	l.logger.Info(args)
}

func (l *Log) Warn(args ...interface{}) {
	l.logger.Warn(args)
}

func (l *Log) Error(args ...interface{}) {
	l.logger.Error(args)
}

func (l *Log) Fatal(args ...interface{}) {
	l.logger.Fatal(args)
	os.Exit(0)
}

func (l *Log) Debugf(template string, args ...interface{}) {
	l.logger.Debug(fmt.Sprintf(template, args...))
}

func (l *Log) Infof(template string, args ...interface{}) {
	l.logger.Info(fmt.Sprintf(template, args...))
}

func (l *Log) Errorf(template string, args ...interface{}) {
	l.logger.Error(fmt.Sprintf(template, args...))
}
