package log

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

type LogLevel string

const (
	DEBUG    LogLevel = "DEBUG"
	INFO     LogLevel = "INFO"
	WARNING  LogLevel = "WARNING"
	ERROR    LogLevel = "ERROR"
	CRITICAL LogLevel = "CRITICAL"
	FATAL    LogLevel = "FATAL"
	PANIC    LogLevel = "PANIC"
)

var LOG_LEVEL_MAP = map[LogLevel]logrus.Level{
	DEBUG:    logrus.DebugLevel,
	INFO:     logrus.InfoLevel,
	WARNING:  logrus.WarnLevel,
	ERROR:    logrus.ErrorLevel,
	CRITICAL: logrus.FatalLevel,
	FATAL:    logrus.FatalLevel,
	PANIC:    logrus.PanicLevel,
}

//Log type represents a logrus wrapper
type Log struct {
	log *logrus.Logger
}

//NewStdOutLogger use from test cases
func NewStdOutLogger(appName string, logLevel int32) *Log {
	l := Log{}
	l.log = logrus.New()
	l.log.SetLevel(logrus.Level(logLevel))
	l.log.Out = os.Stdout
	l.log.Formatter = &BloomLogFormatter{
		TextFormatter: prefixed.TextFormatter{
			FullTimestamp:    false,
			DisableSorting:   true,
			TimestampFormat:  "15:04:05.000",
			DisableTimestamp: true,
		},
		BloomDisbleTs: false,
	}

	return &l
}

//New get new logrus logger with sentry enabled
func New(file *os.File, logLevel int32) *Log {

	var l Log
	l.log = logrus.New()

	fmt.Printf("logLevel: %d\n", logLevel)
	l.log.SetLevel(logrus.Level(logLevel))
	l.log.Formatter = &BloomLogFormatter{
		TextFormatter: prefixed.TextFormatter{
			DisableColors:    true,
			FullTimestamp:    true,
			DisableSorting:   true,
			TimestampFormat:  "2006-01-02 15:04:05.000",
			DisableTimestamp: true,
		},
		BloomDisbleTs: false,
	}

	l.log.SetOutput(file)

	return &l
}

//WithFields pass fields as a map[string]interface{}
func (l *Log) WithFields(fields map[string]interface{}) *logrus.Entry {

	return l.log.WithFields(logrus.Fields(fields))
}

//Info .
func (l *Log) Info(args ...interface{}) {
	l.log.Info(args)
}

//Infof .
func (l *Log) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args)
}

//Infoln .
func (l *Log) Infoln(args ...interface{}) {
	l.log.Infoln(args)
}

//Debug .
func (l *Log) Debug(args ...interface{}) {
	l.log.Debug(args)
}

//Debugf .
func (l *Log) Debugf(format string, args ...interface{}) {
	l.log.Debugf(format, args)
}

//Debugln .
func (l *Log) Debugln(args ...interface{}) {
	l.log.Debugln(args)
}

//Error .
func (l *Log) Error(args ...interface{}) {
	l.log.Error(args)
}

//Errorf .
func (l *Log) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args)
}

//Errorln .
func (l *Log) Errorln(args ...interface{}) {
	l.log.Errorln(args)
}

//Warn .
func (l *Log) Warn(args ...interface{}) {
	l.log.Warn(args)
}

//Warnf .
func (l *Log) Warnf(format string, args ...interface{}) {
	l.log.Warnf(format, args)
}

//Warnln .
func (l *Log) Warnln(args ...interface{}) {
	l.log.Warnln(args)
}

//Fatal .
func (l *Log) Fatal(args ...interface{}) {
	l.log.Fatal(args)
}

//Fatalf .
func (l *Log) Fatalf(format string, args ...interface{}) {
	l.log.Fatalf(format, args)
}

//Fatalln .
func (l *Log) Fatalln(args ...interface{}) {
	l.log.Fatalln(args)
}

//Panic .
func (l *Log) Panic(args ...interface{}) {
	l.log.Panic(args)
}

//Panicf .
func (l *Log) Panicf(format string, args ...interface{}) {
	l.log.Panicf(format, args)
}

//Panicln .
func (l *Log) Panicln(args ...interface{}) {
	l.log.Panicln(args)
}

//Print .
func (l *Log) Print(args ...interface{}) {
	l.log.Debug(args)
}

//Println .
func (l *Log) Println(args ...interface{}) {
	l.log.Debugln(args)
}

//Printf .
func (l *Log) Printf(format string, args ...interface{}) {
	l.log.Printf(format, args)
}

//Wrap TODO: to add func and line number
func Wrap(logger *logrus.Entry) *logrus.Entry {
	// if pc, _, line, ok := runtime.Caller(1); ok {
	// 	fName := runtime.FuncForPC(pc).Name()
	// 	return logger.WithField("line", line).WithField("func", fName)
	// }

	return logger
}

type LoggerWrapperForLogrus struct {
	log *logrus.Entry
}

func NewLoggerFromLogrus(log *logrus.Entry) LoggerWrapperForLogrus {
	return LoggerWrapperForLogrus{log: log}
}
func (l LoggerWrapperForLogrus) Info(msg string, keysAndValues ...interface{}) {
	l.log.Infof(msg, keysAndValues...)
}

func (l LoggerWrapperForLogrus) Error(err error, msg string, keysAndValues ...interface{}) {
	l.log.Errorf(msg, keysAndValues...)
}
