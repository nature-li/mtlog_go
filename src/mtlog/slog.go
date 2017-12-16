package mtlog

import (
	"log/syslog"
	"fmt"
	"time"
)

type SysLog struct {
	writer   *syslog.Writer
	maxCount int
	curCount int
	lastTime int64
}

func newSysLog(maxCount int) *SysLog {
	return &SysLog{
		writer:   nil,
		maxCount: maxCount,
		curCount: 0,
		lastTime: 0,
	}
}

func (o *SysLog) init() bool {
	var err error
	o.writer, err = syslog.New(syslog.LOG_ERR, "mt_log")
	if err != nil {
		fmt.Println("New syslog failed: ", err.Error())
		return false
	}

	return true
}

func (o *SysLog) close() {
	if o.writer != nil {
		err := o.writer.Close()
		if err != nil {
			fmt.Println("close syslog failed: ", err.Error())
		}
	}
}

func (o *SysLog) error(msg string) {
	if o.writer != nil {
		if !o.shouldLog() {
			return
		}

		err := o.writer.Err(msg)
		if err != nil {
			fmt.Println("write syslog failed: ", err.Error())
		}
	}
}

func (o *SysLog) shouldLog() bool {
	curTime := time.Now().Unix()
	seconds := curTime - o.lastTime

	if seconds >= 3600 {
		o.curCount = 0
		o.lastTime = curTime
	}

	if o.curCount >= o.maxCount {
		return false
	}

	o.curCount++
	return true
}
