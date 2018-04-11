package mtlog

import (
	"os"
	"path/filepath"
	"time"
	"io/ioutil"
	"regexp"
	"sort"
	"strings"
)

type fileGroup struct {
	fileDir   string
	fileName  string
	maxSize   int64
	fileCount int
	fileMap   map[Level]*fileInfo
	fileArray []*fileInfo
	timer     *time.Timer
	flag      chan bool
}

func newFileGroup(fileDir string, fileName string, maxSize int64, fileCount int) *fileGroup {
	return &fileGroup{
		fileDir:   fileDir,
		fileName:  fileName,
		maxSize:   maxSize,
		fileCount: fileCount,
		fileMap:   make(map[Level]*fileInfo),
		fileArray: make([]*fileInfo, 0),
		timer:     time.NewTimer(60 * time.Second),
		flag:      make(chan bool, 1),
	}
}

func (o *fileGroup) init() bool {
	// make sure dir exist
	exist := false
	if _, err := os.Stat(o.fileDir); err == nil {
		exist = true
	}

	// create dir
	if !exist {
		err := os.Mkdir(o.fileDir, 0755)
		if err != nil {
			slog.error(err.Error())
			return false
		}
	}

	// open process file
	name := filepath.Join(o.fileDir, o.fileName)
	process := newFileInfo(name+".process.log", o.maxSize)
	if !process.reopen() {
		return false
	}
	o.fileMap[TRACE] = process
	o.fileMap[DEBUG] = process
	o.fileMap[INFO] = process
	o.fileMap[WARN] = process
	o.fileMap[ERROR] = process
	o.fileMap[FATAL] = process
	o.fileArray = append(o.fileArray, process)

	// open report file
	report := newFileInfo(name+".report.log", o.maxSize)
	if !report.reopen() {
		return false
	}
	o.fileMap[REPORT] = report
	o.fileArray = append(o.fileArray, report)

	// start clean thread
	if o.fileCount >= 0 {
		go o.cleanLoop()
	}
	return true
}

func (o *fileGroup) stop() {
	for _, info := range o.fileArray {
		info.close()
	}

	if o.fileCount >= 0 {
		o.flag <- true
	}
}

func (o *fileGroup) write(r *record) bool {
	info, ok := o.fileMap[r.level]
	if !ok {
		slog.error("no file info for level: " + r.level.String())
		return false
	}

	result := info.write(r.level, r.content)
	return result
}

func (o *fileGroup) writeFlushRotate(r *record) bool {
	info, ok := o.fileMap[r.level]
	if !ok {
		slog.error("no file info for level: " + r.level.String())
		return false
	}

	result := info.writeFlushRotate(r.level, r.content)
	return result
}

func (o *fileGroup) flush() {
	for _, info := range o.fileArray {
		info.flush()
	}
}

func (o *fileGroup) rotate() {
	for _, info := range o.fileArray {
		if info.needRotate() {
			info.rotate()
		}
	}
}

func (o *fileGroup) cleanLoop() {
	quit := false

	for !quit {
		select {
		case <-o.timer.C:
			o.clean("process")
			o.clean("report")
			o.timer.Reset(60 * time.Second)

		case <-o.flag:
			quit = true
		}
	}
}

func (o *fileGroup) clean(middle string) {
	files, err := ioutil.ReadDir(o.fileDir)
	if err != nil {
		slog.error(err.Error())
		return
	}

	fileMap := make(map[string]os.FileInfo)
	var reg = regexp.MustCompile("^" + o.fileName + "." + middle + ".log.[0-9]{20}$")
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		if !reg.MatchString(name) {
			continue
		}

		fields := strings.Split(name, ".")
		tail := fields[len(fields)-1]
		if !o.checkTail(tail) {
			continue
		}

		fileMap[name] = file
	}

	o.delete(fileMap)
}

func (o *fileGroup) delete(m map[string]os.FileInfo) {
	more := len(m) - o.fileCount

	if more > 0 {
		var keys []string
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for i := 0; i < more; i++ {
			key := keys[i]
			file := m[key]

			fullName := filepath.Join(o.fileDir, file.Name())
			err := os.Remove(fullName)
			if err != nil {
				slog.error("remove file[" + fullName + "] failed:" + err.Error())
			}
		}
	}
}

func (o *fileGroup) checkTail(tail string) bool {
	timeString := tail[0:14]
	loc, err := time.LoadLocation("Local")
	if err != nil {
		slog.error("LoadLocation error: " + err.Error())
		return false
	}

	when, err := time.ParseInLocation("20060102150405", timeString, loc)
	if err != nil {
		slog.error("ParseInLocation error: " + err.Error())
	}

	toString := when.Format("20060102150405")
	if timeString != toString {
		return false
	}

	return true
}
