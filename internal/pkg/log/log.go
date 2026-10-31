package log

import (
	"io"

	"example.com/go-cli-template/internal/config"
	"github.com/samber/oops"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Init replaces Zap's global logger using the application configuration.
func Init(cfg config.Log, output io.Writer) error {
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return oops.Wrapf(err, "invalid log level")
	}
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	logger := zap.New(zapcore.NewCore(encoder, zapcore.Lock(zapcore.AddSync(output)), level))
	zap.ReplaceGlobals(logger)
	return nil
}
