package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func Info(msg string, fields map[string]interface{}) {
	logFields(fields).Info().Msg(msg)
}

func Error(msg string, err error, fields map[string]interface{}) {
	logFields(fields).Err(err).Msg(msg)
}

func InfoWithRequestID(msg string, requestID string, fields map[string]interface{}) {
	fields["request_id"] = requestID
	Info(msg, fields)
}

func ErrorWithRequestID(msg string, err error, requestID string, fields map[string]interface{}) {
	fields["request_id"] = requestID
	Error(msg, err, fields)
}

func logFields(fields map[string]interface{}) *zerolog.Logger {
	l := log.With()
	for k, v := range fields {
		l = l.Interface(k, v)
	}
	logger := l.Logger()
	return &logger
}
