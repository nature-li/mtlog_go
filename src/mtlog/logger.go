package mtlog

import (
	"fmt"
	"runtime"
	"strconv"
	"os"
)

var (
	queueSize = 100 * 1024
	slog = newSysLog(10)
)

type Logger struct {
	env   Env
	level Level
	sink  *Sink
	pid int
	headerLen int
}

func NewLogger(env Env, level Level, fileDir string, fileName string, maxSize int64, maxCount int) *Logger {
	return &Logger{
		env:   env,
		level: level,
		sink:  newSink(fileDir, fileName, maxSize, maxCount, queueSize),
		pid: os.Getpid(),
		headerLen: len(os.Getenv("GOPATH")),
	}
}

func (o *Logger) Start() bool {
	if !slog.init() {
		return false
	}

	if !o.sink.start() {
		return false
	}

	return true
}

func (o *Logger) Stop() {
	o.sink.stop()
	slog.close()
}

func (o *Logger) SetLevel(level Level) {
	o.level = level
}

func (o *Logger) Trace(content string) {
	if TRACE >= o.level {
		o.log(3, TRACE, "", "normal", content)
	}
}

func (o *Logger) PvTrace(pvId string, content string) {
	if TRACE >= o.level {
		o.log(3, TRACE, pvId, "normal", content)
	}
}

func (o *Logger) PvWordTrace(pvId string, keyword string, content string) {
	if TRACE >= o.level {
		o.log(3, TRACE, pvId, keyword, content)
	}
}

func (o *Logger) Tracef(format string, args ...interface{}) {
	if TRACE >= o.level {
		o.logf(4, TRACE, "", "normal", format, args...)
	}
}

func (o *Logger) PvTracef(pvId string, format string, args ...interface{}) {
	if TRACE >= o.level {
		o.logf(4, TRACE, pvId, "normal", format, args...)
	}
}

func (o *Logger) PvWordTracef(pvId string, keyword string, format string, args ...interface{}) {
	if TRACE >= o.level {
		o.logf(4, TRACE, pvId, keyword, format, args...)
	}
}

func (o *Logger) Debug(content string) {
	if DEBUG >= o.level {
		o.log(3, DEBUG, "", "normal", content)
	}
}

func (o *Logger) PvDebug(pvId string, content string) {
	if DEBUG >= o.level {
		o.log(3, DEBUG, pvId, "normal", content)
	}
}

func (o *Logger) PvWordDebug(pvId string, keyword string, content string) {
	if DEBUG >= o.level {
		o.log(3, DEBUG, pvId, keyword, content)
	}
}

func (o *Logger) Debugf(format string, args ...interface{}) {
	if DEBUG >= o.level {
		o.logf(4, DEBUG, "", "normal", format, args...)
	}
}

func (o *Logger) PvDebugf(pvId string, format string, args ...interface{}) {
	if DEBUG >= o.level {
		o.logf(4, DEBUG, pvId, "normal", format, args...)
	}
}

func (o *Logger) PvWordDebugf(pvId string, keyword string, format string, args ...interface{}) {
	if DEBUG >= o.level {
		o.logf(4, DEBUG, pvId, keyword, format, args...)
	}
}

func (o *Logger) Info(content string) {
	if INFO >= o.level {
		o.log(3, INFO, "", "normal", content)
	}
}

func (o *Logger) PvInfo(pvId string, content string) {
	if INFO >= o.level {
		o.log(3, INFO, pvId, "normal", content)
	}
}

func (o *Logger) PvWordInfo(pvId string, keyword string, content string) {
	if INFO >= o.level {
		o.log(3, INFO, pvId, keyword, content)
	}
}

func (o *Logger) Infof(format string, args ...interface{}) {
	if INFO >= o.level {
		o.logf(4, INFO, "", "normal", format, args...)
	}
}

func (o *Logger) PvInfof(pvId string, format string, args ...interface{}) {
	if INFO >= o.level {
		o.logf(4, INFO, pvId, "normal", format, args...)
	}
}

func (o *Logger) PvWordInfof(pvId string, keyword string, format string, args ...interface{}) {
	if INFO >= o.level {
		o.logf(4, INFO, pvId, keyword, format, args...)
	}
}

func (o *Logger) Warn(content string) {
	if WARN >= o.level {
		o.log(3, WARN, "", "normal", content)
	}
}

func (o *Logger) PvWarn(pvId string, content string) {
	if WARN >= o.level {
		o.log(3, WARN, pvId, "normal", content)
	}
}

func (o *Logger) PvWordWarn(pvId string, keyword string, content string) {
	if WARN >= o.level {
		o.log(3, WARN, pvId, keyword, content)
	}
}

func (o *Logger) Warnf(format string, args ...interface{}) {
	if WARN >= o.level {
		o.logf(4, WARN, "", "normal", format, args...)
	}
}

func (o *Logger) PvWarnf(pvId string, format string, args ...interface{}) {
	if WARN >= o.level {
		o.logf(4, WARN, pvId, "normal", format, args...)
	}
}

func (o *Logger) PvWordWarnf(pvId string, keyword string, format string, args ...interface{}) {
	if WARN >= o.level {
		o.logf(4, WARN, pvId, keyword, format, args...)
	}
}

func (o *Logger) Error(content string) {
	if ERROR >= o.level {
		o.log(3, ERROR, "", "normal", content)
	}
}

func (o *Logger) PvError(pvId string, content string) {
	if ERROR >= o.level {
		o.log(3, ERROR, pvId, "normal", content)
	}
}

func (o *Logger) PvWordError(pvId string, keyword string, content string) {
	if ERROR >= o.level {
		o.log(3, ERROR, pvId, keyword, content)
	}
}

func (o *Logger) Errorf(format string, args ...interface{}) {
	if ERROR >= o.level {
		o.logf(4, ERROR, "", "normal", format, args...)
	}
}

func (o *Logger) PvErrorf(pvId string, format string, args ...interface{}) {
	if ERROR >= o.level {
		o.logf(4, ERROR, pvId, "normal", format, args...)
	}
}

func (o *Logger) PvWordErrorf(pvId string, keyword string, format string, args ...interface{}) {
	if ERROR >= o.level {
		o.logf(4, ERROR, pvId, keyword, format, args...)
	}
}

func (o *Logger) Fatal(content string) {
	if FATAL >= o.level {
		o.log(3, FATAL, "", "normal", content)
	}
}

func (o *Logger) PvFatal(pvId string, content string) {
	if FATAL >= o.level {
		o.log(3, FATAL, pvId, "normal", content)
	}
}

func (o *Logger) PvWordFatal(pvId string, keyword string, content string) {
	if FATAL >= o.level {
		o.log(3, FATAL, pvId, keyword, content)
	}
}

func (o *Logger) Fatalf(format string, args ...interface{}) {
	if FATAL >= o.level {
		o.logf(4, FATAL, "", "normal", format, args...)
	}
}

func (o *Logger) PvFatalf(pvId string, format string, args ...interface{}) {
	if FATAL >= o.level {
		o.logf(4, FATAL, pvId, "normal", format, args...)
	}
}

func (o *Logger) PvWordFatalf(pvId string, keyword string, format string, args ...interface{}) {
	if FATAL >= o.level {
		o.logf(4, FATAL, pvId, keyword, format, args...)
	}
}

func (o *Logger) Report(content string) {
	if REPORT >= o.level {
		o.log(3, REPORT, "", "normal", content)
	}
}

func (o *Logger) PvReport(pvId string, content string) {
	if REPORT >= o.level {
		o.log(3, REPORT, pvId, "normal", content)
	}
}

func (o *Logger) PvWordReport(pvId string, keyword string, content string) {
	if REPORT >= o.level {
		o.log(3, REPORT, pvId, keyword, content)
	}
}

func (o *Logger) Reportf(format string, args ...interface{}) {
	if REPORT >= o.level {
		o.logf(4, REPORT, "", "normal", format, args...)
	}
}

func (o *Logger) PvReportf(pvId string, format string, args ...interface{}) {
	if REPORT >= o.level {
		o.logf(4, REPORT, pvId, "normal", format, args...)
	}
}

func (o *Logger) PvWordReportf(pvId string, keyword string, format string, args ...interface{}) {
	if REPORT >= o.level {
		o.logf(4, REPORT, pvId, keyword, format, args...)
	}
}

func (o *Logger) logf(depth int, level Level, pvId string, keyword string, format string, args ...interface{}) {
	content := fmt.Sprintf(format, args...)
	o.log(depth, level, pvId, keyword, content)
}

// time, level, threadId, position, env, pvId, keyword, content
func (o *Logger) log(depth int, level Level, pvId string, keyword string, content string) {
	// log message
	buf := make([]byte, 0)

	// time
	now := getLogTime()
	buf = append(buf, "["...)
	buf = append(buf, now...)
	buf = append(buf, "]"...)
	buf = append(buf, 0x1e)

	// level
	buf = append(buf, "["...)
	buf = append(buf, level.String()...)
	buf = append(buf, "]"...)
	buf = append(buf, 0x1e)

	// process id
	buf = append(buf, "["...)
	buf = append(buf, strconv.Itoa(o.pid)...)
	buf = append(buf, "]"...)
	buf = append(buf, 0x1e)

	// position
	position := o.getPosition(depth)
	buf = append(buf, "["...)
	buf = append(buf, position...)
	buf = append(buf, "]"...)
	buf = append(buf, 0x1e)

	// env
	buf = append(buf, "["...)
	buf = append(buf, o.env.String()...)
	buf = append(buf, "]"...)
	buf = append(buf, 0x1e)

	// pvId
	buf = append(buf, "["...)
	buf = append(buf, pvId...)
	buf = append(buf, "]"...)
	buf = append(buf, 0x1e)

	// keyword
	buf = append(buf, "["...)
	buf = append(buf, keyword...)
	buf = append(buf, "]"...)
	buf = append(buf, 0x1e)

	// content
	buf = append(buf, "["...)
	buf = append(buf, content...)
	buf = append(buf, "]"...)
	buf = append(buf, 0x1e)

	r := &record{
		level : level,
		content: buf,
	}
	// push to queue
	o.sink.pushBack(r)
}

func (o *Logger) getPosition(depth int) []byte {
	_, fileName, fileLine, ok := runtime.Caller(depth)
	if !ok {
		slog.error("runtime.Caller failed")
		fileName = "???"
		fileLine = 0
	}
	relativeName := fileName[o.headerLen:]

	var buf []byte
	buf = append(buf, relativeName...)
	buf = append(buf, ":"...)
	buf = append(buf, strconv.Itoa(fileLine)...)
	return buf
}