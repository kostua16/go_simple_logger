package logger

import (
	"github.com/op/go-logging"
	"os"
)

//import (
//	"flag"
//	"github.com/op/go-logging"
//)

//var logLevel = *flag.Int("verbose", 2, "Verbose logLevel")

var (
	defaultLevel = logging.WARNING
)

func init() {
	// Example format string. Everything except the message has a custom color
	// which is dependent on the log level. Many fields have a custom output
	// formatting too, eg. the time returns the hour down to the milli second.
	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} | %{module:-8s} | %{longfunc:-60s} | %{level:.4s} | %{id:03x}%{color:reset} ▶ %{message}`,
	)
	// For demo purposes, create two backend for os.Stderr.
	consoleOutput := logging.NewLogBackend(os.Stderr, "", 0)
	consoleOutputFormatted := logging.NewBackendFormatter(consoleOutput, format)

	logging.SetBackend(consoleOutputFormatted)
}

func SetLevel(level logging.Level) {
	defaultLevel = level
}

func SetLevelStr(level string) {
	parsed, parseErr := logging.LogLevel(level)
	if parseErr != nil {
		parsed = logging.WARNING
	}
	SetLevel(parsed)
}

func SetModuleLevel(level logging.Level, name string) {
	logging.SetLevel(level, name)
}

func SetModuleLevelStr(level string, name string) {
	parsed, parseErr := logging.LogLevel(level)
	if parseErr != nil {
		parsed = logging.WARNING
	}
	SetModuleLevel(parsed, name)
}

func CreateLogger(name string) *logging.Logger {

	var log = logging.MustGetLogger(name)
	logging.SetLevel(defaultLevel, name)
	//switch logLevel {
	//case 1:
	//	logging.SetLevel(logging.ERROR, name)
	//	break
	//case 2:
	//	logging.SetLevel(logging.WARNING, name)
	//	break
	//case 3:
	//	logging.SetLevel(logging.INFO, name)
	//	break
	//case 4:
	//	logging.SetLevel(logging.DEBUG, name)
	//	break
	//default:
	//	logging.SetLevel(logging.WARNING, name)
	//	break
	//}
	return log
}

//func init() {
//	consoleWriter := zerolog.ConsoleWriter{Out: os.Stderr}
//	//multi := zerolog.MultiLevelWriter(consoleWriter, os.Stderr)
//	logger := zerolog.New(consoleWriter).With().Timestamp().Caller().Logger()
//	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
//	zerolog.SetGlobalLevel(zerolog.InfoLevel)
//	log.Logger = logger
//}
//
//func SetLevel(level zerolog.Level) {
//	zerolog.SetGlobalLevel(level)
//}
//
//func CreateLogger(name string) zerolog.Logger {
//	return log.With().Str("component", name).Logger()
//}
