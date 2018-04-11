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

func (o *fileInfo) reopen() bool {
	if !o.closed {
		return true
	}

	f, err := os.OpenFile(o.name, os.O_CREATE | os.O_APPEND | os.O_RDWR, 0666)
	if err != nil {
		slog.error(err.Error())
		return false
	}

	if o.w == nil {
		o.w = bufio.NewWriter(f)
	} else {
		o.w.Reset(f)
	}

	stat, err := f.Stat()
	if err != nil {
		slog.error(err.Error())
		return false
	}

	o.f = f
	o.closed = false
	o.needFlush = false
	o.curLen = stat.Size()
	o.lastRotate = stat.ModTime()

	return true
}

func (o *fileInfo) reset() {
	o.f = nil
	o.closed = true
	o.needFlush = false
	o.curLen = 0
	o.lastRotate = time.Now()
}

func (o *fileInfo) close() {
	if o.closed {
		return
	}

	err := o.w.Flush()
	if err != nil {
		slog.error(err.Error())
	}

	err = o.f.Close()
	if err != nil {
		slog.error(err.Error())
	}

	o.reset()
}

func (o *fileInfo) delete() {
	err := os.Remove(o.name)
	if err != nil {
		line := fmt.Sprintf("remove %v failed: %v", o.name, err.Error())
		slog.error(line)
	}
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

	_, err := o.w.Write(content)
	if err != nil {
		slog.error("write file error for level[" + level.String() + "]" + ": " + err.Error())
		return false
	}

	o.curLen += int64(len(content))
	o.needFlush = true

	if o.curLen >= o.maxLen {
		o.rotate()
	}
	return true
}

func (o *fileInfo) writeFlushRotate(level Level, content []byte) bool {
	if o.closed {
		slog.error("file has been closed for level: " + level.String())
		return false
	}

	_, err := o.w.Write(content)
	if err != nil {
		slog.error("write file error for level[" + level.String() + "]" + ": " + err.Error())
		return false
	}

	err = o.w.Flush()
	if err != nil {
		slog.error("flush file error for level[" + level.String() + "]" + ": " + err.Error())
		return false
	}

	o.curLen += int64(len(content))

	if o.curLen >= o.maxLen {
		o.rotate()
	}
	return true
}

func (o *fileInfo) flush() {
	if !o.needFlush {
		return
	}

	err := o.w.Flush()
	if err != nil {
		slog.error(err.Error())
	}
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
	if o.curLen > 0 {
		o.rename()
	} else {
		o.delete()
	}
	o.reopen()
	return true
}