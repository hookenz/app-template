package logging

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type MoneyLogger struct {
	zerolog.Logger
}

var Logger MoneyLogger

func NewLogger() MoneyLogger {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}

	output.FormatErrFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s: ", i)
	}

	zerolog := zerolog.New(output).With().Timestamp().Logger()
	Logger = MoneyLogger{zerolog}
	return Logger
}

func (l *MoneyLogger) LogInfo() *zerolog.Event {
	return l.Logger.Info()
}

func (l *MoneyLogger) LogError() *zerolog.Event {
	return l.Logger.Error()
}

func (l *MoneyLogger) LogDebug() *zerolog.Event {
	return l.Logger.Debug()
}

func (l *MoneyLogger) LogWarn() *zerolog.Event {
	return l.Logger.Warn()
}

func (l *MoneyLogger) LogFatal() *zerolog.Event {
	return l.Logger.Fatal()
}
