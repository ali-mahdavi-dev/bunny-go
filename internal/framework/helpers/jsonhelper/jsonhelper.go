package jsonhelper

import (
	"encoding/json"

	"github.com/ali-mahdavi-dev/bunny-go/config"
	"github.com/ali-mahdavi-dev/bunny-go/internal/framework/infrastructure/logging"
)

var loggging = logging.NewLogger(config.GetConfig())

func Encode[T any](t T) []byte {
	b, err := json.Marshal(t)
	if err != nil {
		loggging.Error(logging.IO, logging.CanNotMarshal, "couldn't encode the variable", map[logging.ExtraKey]interface{}{
			logging.JsonMarshalKey: t,
		})
	}
	return b
}

func Decode[T any](b []byte) T {
	var t T
	err := json.Unmarshal(b, &t)
	if err != nil {
		loggging.Error(logging.IO, logging.CanNotMarshal, "couldn't decode the variable", map[logging.ExtraKey]interface{}{
			logging.JsonMarshalKey:   t,
			logging.JsonMarshalValue: b,
		})
	}
	return t
}
