package jsonhelper

import (
	"encoding/json"

	"shikposh-backend/config"
	"shikposh-backend/pkg/framework/infrastructure/logging"
)

var loggging logging.Logger

func init() {
	cfg := config.GetConfig()
	loggerConfig := logging.LoggerConfig{
		Type:   logging.LoggerTypeZerolog,
		Level:  logging.LogLevel(cfg.Logger.Level),
		Format: logging.LogFormatJSON,
	}
	var err error
	loggging, err = logging.NewLogger(loggerConfig)
	if err != nil {
		panic(err)
	}
}

func Encode[T any](t T) []byte {
	b, err := json.Marshal(t)
	if err != nil {
		logging.Error("JSON encoding failed").
			WithString("operation", "encode").
			WithAny("variable", t).
			Log()
	}
	return b
}

func Decode[T any](b []byte) T {
	var t T
	err := json.Unmarshal(b, &t)
	if err != nil {
		logging.Error("JSON decoding failed").
			WithString("operation", "decode").
			WithAny("variable", t).
			WithAny("bytes", b).
			Log()
	}
	return t
}
