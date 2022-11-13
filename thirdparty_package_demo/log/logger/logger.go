package logger

import (
	"flag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
)

const (
	DefaultTimeFormat = "2006-01-02 15:04:05"
	DefaultLogFile    = "./app.log"
	DefaultMaxSize    = 1 << 30
	DefaultMaxAge     = 7
	DefaultMaxBackups = 3

	ConsoleEncoder = "console"
	JsonEncoder    = "json"

	StdoutWriter = "stdout"
	FileWriter   = "file"
)

var (
	G        = GetLogger
	instance *Logger
)

var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

type Options struct {
	Level      string
	Writer     string
	Encoder    string
	TimeFormat string
	Caller     bool
	Stacktrace bool
	LogFile    string
	MaxSize    int
	MaxAge     int
	MaxBackups int
}

type Logger struct {
	sugared      *zap.SugaredLogger
	globalFields []interface{}
	tempFields   []interface{}
}

type opt func(config *zapcore.EncoderConfig)

func (o *Options) BindFlags(fs *flag.FlagSet) {
	fs.StringVar(&o.Level, "log.level", o.Level, "Log level. Can be one of 'debug', 'info', 'warn', 'error', 'dpanic', 'panic', 'fatal'. Default 'info'.")
	fs.StringVar(&o.Writer, "log.writer", o.Writer, "Log output writer. Can be one of 'stdout', 'file'. Default both.")
	fs.StringVar(&o.Encoder, "log.encoder", o.Encoder, "Log encoder. Can be one of 'console', 'json'. Default 'json'.")
	fs.StringVar(&o.TimeFormat, "log.timeformat", o.TimeFormat, "Log time format. Default '2006-01-02 15:04:05'.")
	fs.BoolVar(&o.Caller, "log.caller", o.Caller, "Log caller. Default false.")
	fs.BoolVar(&o.Stacktrace, "log.stacktrace", o.Stacktrace, "Log stacktrace. Default false.")
	fs.StringVar(&o.LogFile, "log.file", o.LogFile, "Log file path. Default './app.log'.")
	fs.IntVar(&o.MaxSize, "log.filesize", o.MaxSize, "Log file max size. Default ''.")
	fs.IntVar(&o.MaxAge, "log.fileage", o.MaxAge, "Log file max age. Default 7.")
	fs.IntVar(&o.MaxBackups, "log.filebackups", o.MaxBackups, "Log file max backups. Default 3.")
}

func (o *Options) setDefault() {
	if o.Level == "" {
		o.Level = "info"
	}

	if o.Encoder == "" {
		o.Encoder = "json"
	}

	if o.TimeFormat == "" {
		o.TimeFormat = DefaultTimeFormat
	}

	if o.LogFile == "" {
		o.LogFile = DefaultLogFile
	}

	if o.MaxSize == 0 {
		o.MaxSize = DefaultMaxSize
	}

	if o.MaxAge == 0 {
		o.MaxAge = DefaultMaxAge
	}

	if o.MaxBackups == 0 {
		o.MaxBackups = DefaultMaxBackups
	}
}

func Init(o *Options, opts ...opt) {
	o.setDefault()

	var writer zapcore.WriteSyncer
	if strings.ToLower(o.Writer) == StdoutWriter {
		writer = zapcore.AddSync(os.Stdout)
	} else if strings.ToLower(o.Writer) == FileWriter {
		writer = zapcore.AddSync(&lumberjack.Logger{
			Filename:   o.LogFile,
			MaxSize:    o.MaxSize,
			MaxAge:     o.MaxAge,
			MaxBackups: o.MaxBackups,
			LocalTime:  true,
			Compress:   false,
		})
	} else {
		writer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&lumberjack.Logger{
			Filename:   o.LogFile,
			MaxSize:    o.MaxSize,
			MaxAge:     o.MaxAge,
			MaxBackups: o.MaxBackups,
			LocalTime:  true,
			Compress:   false,
		}))
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(o.TimeFormat)
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.MessageKey = "message"
	encoderConfig.LevelKey = "err_level"

	if opts != nil {
		for _, f := range opts {
			f(&encoderConfig)
		}
	}

	var encoder interface{}
	if strings.ToLower(o.Encoder) == ConsoleEncoder {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}
	core := zapcore.NewCore(encoder.(zapcore.Encoder), writer, zap.NewAtomicLevelAt(getLogLevel(strings.ToLower(o.Level))))

	var zos []zap.Option
	if o.Caller {
		zos = append(zos, zap.AddCaller())
	}
	if o.Stacktrace {
		zos = append(zos, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	l := zap.New(core, zos...).Sugar()

	instance = &Logger{
		sugared:      l,
		globalFields: nil,
		tempFields:   nil,
	}
}

func getLogLevel(lvl string) zapcore.Level {
	if level, ok := levelMap[lvl]; ok {
		return level
	}
	return zapcore.InfoLevel
}

func GetLogger() *Logger {
	if instance == nil {
		panic("default logger not initialized")
	}

	return instance
}

func (l *Logger) AddGlobalFields(args ...interface{}) {
	l.globalFields = args
	l.sugared = l.sugared.With(args...)
}

func (l *Logger) AddTempFields(args ...interface{}) {
	if l.tempFields == nil {
		l.tempFields = []interface{}{}
	}
	l.tempFields = append(l.tempFields, args...)
}

func (l *Logger) ResetTempFields() {
	l.tempFields = nil
}

func (l *Logger) With(keysAndValues ...interface{}) *zap.SugaredLogger {
	return l.sugared.With(keysAndValues...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.sugared.With(l.tempFields...).Debug(args...)
}

func (l *Logger) Debugf(template string, args ...interface{}) {
	l.sugared.With(l.tempFields...).Debugf(template, args...)
}

func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.sugared.With(l.tempFields...).Debugw(msg, keysAndValues...)
}

func (l *Logger) Info(args ...interface{}) {
	l.sugared.With(l.tempFields...).Info(args...)
}

func (l *Logger) Infof(template string, args ...interface{}) {
	l.sugared.With(l.tempFields...).Infof(template, args...)
}

func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	l.sugared.With(l.tempFields...).Infow(msg, keysAndValues...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.sugared.With(l.tempFields...).Warn(args...)
}

func (l *Logger) Warnf(template string, args ...interface{}) {
	l.sugared.With(l.tempFields...).Warnf(template, args...)
}

func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.sugared.With(l.tempFields...).Warnw(msg, keysAndValues...)
}

func (l *Logger) Error(args ...interface{}) {
	l.sugared.With(l.tempFields...).Error(args...)
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	l.sugared.With(l.tempFields...).Errorf(template, args...)
}

func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.sugared.With(l.tempFields...).Errorw(msg, keysAndValues...)
}

func (l *Logger) DPanic(args ...interface{}) {
	l.sugared.With(l.tempFields...).DPanic(args...)
}

func (l *Logger) DPanicf(template string, args ...interface{}) {
	l.sugared.With(l.tempFields...).DPanicf(template, args...)
}

func (l *Logger) DPanicw(msg string, keysAndValues ...interface{}) {
	l.sugared.With(l.tempFields...).DPanicw(msg, keysAndValues...)
}

func (l *Logger) Panic(args ...interface{}) {
	l.sugared.With(l.tempFields...).Panic(args...)
}

func (l *Logger) Panicf(template string, args ...interface{}) {
	l.sugared.With(l.tempFields...).Panicf(template, args...)
}

func (l *Logger) Panicw(msg string, keysAndValues ...interface{}) {
	l.sugared.With(l.tempFields...).Panicw(msg, keysAndValues...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.sugared.With(l.tempFields...).Fatal(args...)
}

func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.sugared.With(l.tempFields...).Fatalf(template, args...)
}

func (l *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.sugared.With(l.tempFields...).Fatalw(msg, keysAndValues...)
}
