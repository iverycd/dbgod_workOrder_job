package loggerzap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var (
	globalLogger *zap.Logger
)

const customTimeFormat = "2006-01-02 15:04:05.000"

// 初始化日志系统v1
//func InitV1(env string, cfg *LogConfig) error {
//	var core zapcore.Core
//	// 设置编码器
//	encoderConfig := zap.NewProductionEncoderConfig()
//	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
//	// 不同级别的日志用不同颜色
//	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
//	//encoder := zapcore.NewJSONEncoder(encoderConfig) // json格式存储
//	encoder := zapcore.NewConsoleEncoder(encoderConfig) // 终端输出转储到日志文件
//	// 生产环境配置
//	if env == "production" {
//		// 日志文件切割配置
//		lumberJackLogger := &lumberjack.Logger{
//			Filename:   cfg.FilePath,
//			MaxSize:    cfg.MaxSize,
//			MaxBackups: cfg.MaxBackups,
//			MaxAge:     cfg.MaxAge,
//			Compress:   cfg.Compress,
//			LocalTime:  true,
//		}
//
//		// 生产环境使用文件+错误级别过滤
//		core = zapcore.NewCore(
//			encoder,
//			zapcore.AddSync(lumberJackLogger),
//			getZapLevel(cfg.Level),
//		)
//
//	} else {
//		// 开发环境使用控制台+彩色输出
//		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
//		core = zapcore.NewCore(
//			consoleEncoder,
//			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
//			zapcore.DebugLevel,
//		)
//	}
//	// 创建Logger
//	globalLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
//
//	// 替换zap全局Logger
//	zap.ReplaceGlobals(globalLogger)
//	return nil
//}

// Init 初始化日志系统v2 同时在终端与文件输出日志，使用平面文件格式，彩色输出
func Init(env string, cfg *LogConfig) error {
	var core zapcore.Core
	// 自定义时间编码器
	customTimeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(customTimeFormat))
	}
	// 生产环境配置
	if env == "production" {
		// 文件日志输出配置
		fileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.FilePath,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
			LocalTime:  true,
		})
		// 控制台输出配置
		consoleWriteSyncer := zapcore.AddSync(os.Stdout)

		// 共享的编码器配置基础
		baseEncoderConfig := zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     customTimeEncoder, //zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
		// 文件编码器 - 平面文件格式
		fileEncoder := zapcore.NewConsoleEncoder(baseEncoderConfig)

		// 控制台编码器 - 带颜色的控制台格式
		consoleEncoderConfig := baseEncoderConfig
		consoleEncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder // 终端带颜色
		consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)

		// 创建核心 - 确保文件输出级别足够低
		fileCore := zapcore.NewCore(
			//zapcore.NewJSONEncoder(encoderConfig),// json格式
			fileEncoder, // 平面文件格式编码器
			fileWriteSyncer,
			zapcore.InfoLevel, // 文件记录所有级别日志
		)

		consoleCore := zapcore.NewCore(
			//zapcore.NewConsoleEncoder(encoderConfig),
			consoleEncoder, // 控制台格式编码器
			consoleWriteSyncer,
			zapcore.InfoLevel, // 终端只记录Info及以上级别
		)

		// 合并多个输出
		//multiWriteSyncer := zapcore.NewMultiWriteSyncer(fileWriteSyncer, consoleWriteSyncer)
		// 创建核心 生产环境使用文件+错误级别过滤
		//core = zapcore.NewCore(
		//	zapcore.NewJSONEncoder(encoderConfig), // 编码器
		//	multiWriteSyncer,                      // 输出目标
		//	getZapLevel(cfg.Level),
		//)

		// 合并多个核心
		core = zapcore.NewTee(fileCore, consoleCore)

	} else {
		// 开发环境使用控制台+彩色输出
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewCore(
			consoleEncoder,
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
			zapcore.DebugLevel,
		)
	}
	// 创建Logger
	globalLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

	// 替换zap全局Logger
	zap.ReplaceGlobals(globalLogger)
	return nil
}

// 获取日志级别
func getZapLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// 获取全局Logger
func L() *zap.Logger {
	return globalLogger
}

// 安全关闭
func Sync() error {
	return globalLogger.Sync()
}
