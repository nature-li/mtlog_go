package mtlog

import (
	"time"
)

type Sink struct {
	ch    chan interface{}
	flag  chan bool
	done  chan bool
	group *fileGroup
	timer *time.Timer
	async bool
}

func newSink(fileDir string, fileName string, maxSize int64, fileCount int, queueSize int) *Sink {
	return &Sink{
		ch:    make(chan interface{}, queueSize),
		flag:  make(chan bool, 1),
		done:  make(chan bool, 0),
		group: newFileGroup(fileDir, fileName, maxSize, fileCount),
		timer: time.NewTimer(time.Second * 5),
		async: true,
	}
}

func (o *Sink) setAsync(async bool) {
	o.async = async
}

func (o *Sink) start() bool {
	if !o.group.init() {
		return false
	}

	if o.async {
		go o.consume()
	}
	return true
}

func (o *Sink) stop() {
	if o.async {
		o.flag <- true
		<-o.done
	}

	o.group.stop()
}

func (o *Sink) pushBack(v interface{}) {
	if o.async {
		// push item to queue
		o.ch <- v
	} else {
		// write item to disk
		r := v.(*record)
		o.group.writeFlushRotate(r)
	}
}

func (o *Sink) handleQueue(v interface{}) {
	if v != nil {
		r := v.(*record)
		o.group.write(r)
	}

	for len(o.ch) != 0 {
		v := <-o.ch
		r := v.(*record)
		o.group.write(r)
	}

	o.group.flush()
}

func (o *Sink) consume() {
	quit := false

	for !quit {
		select {
		case v := <-o.ch:
			o.handleQueue(v)

		case <-o.timer.C:
			o.group.rotate()
			o.timer.Reset(5 * time.Second)

		case <-o.flag:
			quit = true
		}
	}

	o.handleQueue(nil)
	o.done <- true
}
