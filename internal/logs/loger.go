
package logs

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger() {
    config := zap.NewProductionConfig()
    config.OutputPaths = []string{
        "logs/app.log", 
        "stdout",      
    }

    config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder 
    logger, err := config.Build()
    if err != nil {
        panic(err)
    }

    Logger = logger
}

// این تابع برای بستن لاگر بعد از استفاده است
func Sync() {
    if Logger != nil {
        _ = Logger.Sync()
    }
}
