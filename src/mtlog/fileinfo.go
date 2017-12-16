package mtlog

import (
	"os"
	"bufio"
	"time"
	"fmt"
)

type fileInfo struct {
	name string
	maxLen int64
	f *os.File
	w *bufio.Writer
	closed bool
	needFlush bool
	curLen int64
	lastRotate time.Time
}

func newFileInfo(name string, maxLen int64) *fileInfo {
	return &fileInfo{
		name: name,
		maxLen: maxLen,
		f: nil,
		w: nil,
		closed: true,
		needFlush: false,
		curLen: 0,
		lastRotate: time.Now(),
	}
}

func (o *fileInfo) open() bool {
	if !o.closed {
		return true
	}

	f, err := os.OpenFile(o.name, os.O_CREATE | os.O_APPEND | os.O_RDWR, 0666)
	if err != nil {
		slog.error(err.Error())
		return false
	}

	w := bufio.NewWriter(f)

	stat, err := f.Stat()
	if err != nil {
		slog.error(err.Error())
		return false
	}

	o.f = f
	o.w = w
	o.closed = false
	o.needFlush = false
	o.curLen = stat.Size()
	o.lastRotate = stat.ModTime()

	return true
}

func (o *fileInfo) reset() {
	o.f = nil
	o.w = nil
	o.closed = true
	o.needFlush = false
	o.curLen = 0
	o.lastRotate = time.Now()
}

func (o *fileInfo) close() {
	if o.closed {
		return
	}

	o.w.Flush()
	o.f.Close()

	o.reset()
}

func (o *fileInfo) rename() {
	newName := o.name + "." + string(getFileTime())
	err := os.Rename(o.name, newName)
	if err != nil {
		line := fmt.Sprintf("rename %v to %v failed: %v", o.name, newName, err.Error())
		slog.error(line)
	}
}

func (o *fileInfo) write(level Level, content []byte) bool {
	if o.closed {
		slog.error("file has been closed for level: " + level.String())
		return false
	}

	_, err := fmt.Fprintln(o.w, string(content))
	if err != nil {
		slog.error("write file error for level[" + level.String() + "]" + ": " + err.Error())
		return false
	}

	o.curLen += int64(len(content))
	o.needFlush = true
	return true
}

func (o *fileInfo) flush() {
	if !o.needFlush {
		return
	}

	o.w.Flush()
	o.needFlush = false
}

func (o *fileInfo) needRotate() bool {
	if o.curLen >= o.maxLen {
		return true
	}

	_, _, lastDay := o.lastRotate.Date()
	_, _, thisDay := time.Now().Date()
	if lastDay != thisDay {
		return true
	}

	return false
}

func (o *fileInfo) rotate() bool {
	o.close()
	o.rename()
	o.open()
	return true
}