package log

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func SetupLogger() (*zap.Logger, error) {
	lumlog.Filename = fmt.Sprintf("./logs/%s_profile-builder-backend.log", time.Now().Format("2006-01-02_15-04-05"))
	file, err := os.Create(lumlog.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %w", err)
	}
	file.Close()

	logger, err := zap.NewProduction(zap.Hooks(lumberjackZapHook))
	if err != nil {
		return nil, fmt.Errorf("failed to create zap logger: %w", err)
	}
	zap.ReplaceGlobals(logger)

	logger.Info("Logger successfully initialized")
	return logger, nil
}

func lumberjackZapHook(e zapcore.Entry) error {
	_, err := lumlog.Write([]byte(fmt.Sprintf("%+v\n", e)))
	if err != nil {
		zap.S().Error("Error writing to log file:", err)
	}
	return err
}
