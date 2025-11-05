package jsonhelper

import (
	"encoding/json"

	"shikposh-backend/pkg/framework/infrastructure/logging"
)

func Encode[T any](t T) []byte {
	b, err := json.Marshal(t)
	if err != nil {
		// Use package-level logging functions
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
		// Use package-level logging functions
		logging.Error("JSON decoding failed").
			WithString("operation", "decode").
			WithAny("variable", t).
			WithAny("bytes", b).
			Log()
	}
	return t
}
