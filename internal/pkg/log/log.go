package log

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
)

// Level 是从logrus借鉴的日志登记，当前日志底层都是借助logrus来做的。
type Level uint32

// Level 的枚举值。
const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel Level = iota
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	// 因为被包裹了一层，默认方法取的当前行号错了一层。关闭默认的打印函数行数的函数，我来手工给写上。
	logrus.SetReportCaller(false)
	logrus.SetLevel(logrus.DebugLevel)
}

// SetLevel sets the log implementation level.
func SetLevel(level Level) {
	logrus.SetLevel(logrus.Level(level))
}

func SetFormatter(formatter logrus.Formatter) {
	logrus.SetFormatter(formatter)
}

// Enable checks if the log level of the log implementation is greater than the level param
func Enable(level Level) bool {
	return logrus.IsLevelEnabled(logrus.Level(level))
}

func buildLog(ctx context.Context) *logrus.Entry {
	return logrus.WithField("file", callerLine(skipForThisFunc)).
		WithField("func", funcName(skipForThisFunc))
}

// Trace logs a message at level Trace on the log implementation.
func Trace(ctx context.Context, args ...interface{}) {
	buildLog(ctx).Trace(args...)
}

// Debug logs a message at level Debug on the log implementation.
func Debug(ctx context.Context, args ...interface{}) {
	buildLog(ctx).Debug(args...)
}

// Print logs a message at level Info on the log implementation.
func Print(ctx context.Context, args ...interface{}) {
	buildLog(ctx).Print(args...)
}

// Info logs a message at level Info on the log implementation.
func Info(ctx context.Context, args ...interface{}) {
	buildLog(ctx).Info(args...)
}

// Warn logs a message at level Warn on the log implementation.
func Warn(ctx context.Context, args ...interface{}) {
	buildLog(ctx).Warn(args...)
}

// Warning logs a message at level Warn on the log implementation.
func Warning(ctx context.Context, args ...interface{}) {
	buildLog(ctx).Warning(args...)
}

// Error logs a message at level Error on the log implementation.
func Error(ctx context.Context, args ...interface{}) {
	buildLog(ctx).Error(args...)
}

// Panic logs a message at level Panic on the log implementation.
func Panic(ctx context.Context, args ...interface{}) {
	buildLog(ctx).Panic(args...)
}

// Fatal logs a message at level Fatal on the log implementation then the process will exit with status set to 1.
func Fatal(ctx context.Context, args ...interface{}) {
	buildLog(ctx).Fatal(args...)
}

// Tracef logs a message at level Trace on the log implementation.
func Tracef(ctx context.Context, format string, args ...interface{}) {
	buildLog(ctx).Tracef(format, args...)
}

// Debugf logs a message at level Debug on the log implementation.
func Debugf(ctx context.Context, format string, args ...interface{}) {
	buildLog(ctx).Debugf(format, args...)
}

// Printf logs a message at level Info on the log implementation.
func Printf(ctx context.Context, format string, args ...interface{}) {
	buildLog(ctx).Printf(format, args...)
}

// Infof logs a message at level Info on the log implementation.
func Infof(ctx context.Context, format string, args ...interface{}) {
	buildLog(ctx).Infof(format, args...)
}

// Warnf logs a message at level Warn on the log implementation.
func Warnf(ctx context.Context, format string, args ...interface{}) {
	buildLog(ctx).Warnf(format, args...)
}

// Warningf logs a message at level Warn on the log implementation.
func Warningf(ctx context.Context, format string, args ...interface{}) {
	buildLog(ctx).Warningf(format, args...)
}

// Errorf logs a message at level Error on the log implementation.
func Errorf(ctx context.Context, format string, args ...interface{}) {
	buildLog(ctx).Errorf(format, args...)
}

// Panicf logs a message at level Panic on the log implementation.
func Panicf(ctx context.Context, format string, args ...interface{}) {
	buildLog(ctx).Panicf(format, args...)
}

// Fatalf logs a message at level Fatal on the log implementation then the process will exit with status set to 1.
func Fatalf(ctx context.Context, format string, args ...interface{}) {
	buildLog(ctx).Fatalf(format, args...)
}

// Traceln logs a message at level Trace on the log implementation.
func Traceln(ctx context.Context, args ...interface{}) {
	buildLog(ctx).Traceln(args...)
}

// Debugln logs a message at level Debug on the log implementation.
func Debugln(ctx context.Context, args ...interface{}) {
	buildLog(ctx).Debugln(args...)
}

// Println logs a message at level Info on the log implementation.
func Println(ctx context.Context, args ...interface{}) {
	buildLog(ctx).Println(args...)
}

// Infoln logs a message at level Info on the log implementation.
func Infoln(ctx context.Context, args ...interface{}) {
	buildLog(ctx).Infoln(args...)
}

// Warnln logs a message at level Warn on the log implementation.
func Warnln(ctx context.Context, args ...interface{}) {
	buildLog(ctx).Warnln(args...)
}

// Warningln logs a message at level Warn on the log implementation.
func Warningln(ctx context.Context, args ...interface{}) {
	buildLog(ctx).Warningln(args...)
}

// Errorln logs a message at level Error on the log implementation.
func Errorln(ctx context.Context, args ...interface{}) {
	buildLog(ctx).Errorln(args...)
}

// Panicln logs a message at level Panic on the log implementation.
func Panicln(ctx context.Context, args ...interface{}) {
	buildLog(ctx).Panicln(args...)
}

// Fatalln logs a message at level Fatal on the log implementation then the process will exit with status set to 1.
func Fatalln(ctx context.Context, args ...interface{}) {
	buildLog(ctx).Fatalln(args...)
}
