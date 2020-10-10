package main

import (
	"log"
	"os"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	level      = zapcore.DebugLevel
	fileLogger *zap.Logger
	fileSugar  *zap.SugaredLogger

	LogFieldSeperator = "|"

	prodLogger *zap.Logger
)

var zapConfig = zapcore.EncoderConfig{
	MessageKey: "msg",
	LevelKey:   "level",
	// Keys can be anything except the empty string.
	TimeKey:          "time",
	NameKey:          "nn",
	CallerKey:        "caller",
	FunctionKey:      "fk",
	StacktraceKey:    "st", // "" 不打堆栈, 非"" 就打印堆栈
	LineEnding:       zapcore.DefaultLineEnding,
	EncodeLevel:      LevelEncoder,
	EncodeTime:       TimeEncoder,
	EncodeDuration:   zapcore.NanosDurationEncoder,
	EncodeCaller:     CallerEncoder,
	EncodeName:       zapcore.FullNameEncoder,
	ConsoleSeparator: "",
}

// TimeEncoder 时间格式
func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000000"))
	enc.AppendString(LogFieldSeperator)
	enc.AppendString(getLocalIP())
	enc.AppendString(LogFieldSeperator)
}

func getLocalIP() string {
	return "127.0.0.1"
}

// CallerEncoder 方法
func CallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	zapcore.ShortCallerEncoder(caller, enc)
	enc.AppendString(LogFieldSeperator)
}

// LevelEncoder 级别格式
func LevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(l.String())
	enc.AppendString(LogFieldSeperator)
}

func Init() {
	var filename = "log/run/csg_run.log"

	lj := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    20, // MB
		MaxAge:     30, // 天
		MaxBackups: 10, // 1表示只保留一个log历史文件, 其他的都压缩, 0表示保留所有log历史, 不进行压缩
		LocalTime:  true,
		Compress:   false, // true 表示会压缩成gz, 而不是打包成zip
	}
	log.Println(lj)

	fileCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zapConfig),
		// zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lj)),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		// zapcore.NewMultiWriteSyncer(zapcore.AddSync(lj)),
		level,
	)

	fileCore.With([]zap.Field{zap.String("vvv", "kkk")})
	fileLogger = zap.New(fileCore, zap.AddCaller(), zap.AddCallerSkip(0), zap.AddStacktrace(zap.DebugLevel))
	fileLogger.WithOptions(zap.Fields(zap.String("service", "csg")))
	// defer fileLogger.Sync()
	fileLogger.Core().Sync()
	fileSugar = fileLogger.Sugar()

	log.Println("[ZapLog] Init success, level: ", level)

	prodLogger, _ = zap.NewProduction(zap.AddStacktrace(zapcore.WarnLevel))
}

func main() {
	Init()

	prodLogger.Error(":cc")

	fileSugar.Info("sxxs")

	fileSugar.DPanicf("cc:")
	fileSugar.Infow("fgb")
	fileLogger.Info("sss")
	fileLogger.Error("xx")
	fileSugar.Fatal("eex")
}
