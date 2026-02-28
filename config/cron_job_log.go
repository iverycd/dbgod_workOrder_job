package config

import "log"

// CustomLogger 自定义日志
type CustomLogger struct{}

func (l *CustomLogger) Debugf(format string, args ...any) { log.Printf("[DEBUG] "+format, args...) }
func (l *CustomLogger) Infof(format string, args ...any)  { log.Printf("[INFO] "+format, args...) }
func (l *CustomLogger) Warnf(format string, args ...any)  { log.Printf("[WARN] "+format, args...) }
func (l *CustomLogger) Errorf(format string, args ...any) { log.Printf("[ERROR] "+format, args...) }
