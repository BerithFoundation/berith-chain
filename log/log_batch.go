package log

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

type BerithNetHandler func(string)

type BerithLogBatch struct {
	ch           chan *Record
	file         *os.File
	time         time.Time
	rotatePeriod time.Duration
	logdir       string
	format       Format
	StopCh       chan interface{}
	cnt          int
	buffer       string
	handler      BerithNetHandler
}
type logPost struct {
	Enode      string `json:"enode"`
	Berithbase string `json:"berithbase"`
	Logs       string `json:"logs"`
}

func NewBerithLogBatch(ch chan *Record, logdir string, rotatePeriod time.Duration, format Format) *BerithLogBatch {
	return &BerithLogBatch{
		ch:           ch,
		logdir:       logdir,
		rotatePeriod: rotatePeriod,
		format:       format,
		StopCh:       make(chan interface{}),
		handler:      func(string) {},
		buffer:       "",
	}
}

func (b *BerithLogBatch) SetHandler(handler BerithNetHandler) {
	b.handler = handler
}

func (b *BerithLogBatch) Loop() {
	for {
		select {
		case record := <-b.ch:
			b.cnt++
			if b.file == nil || time.Now().Sub(b.time) >= b.rotatePeriod {

				if err := os.MkdirAll(b.logdir, 0700); err != nil {
					continue
				}
				now := time.Now()
				logpath := filepath.Join(b.logdir, strings.Replace(now.Format("060102150405.00"), ".", "", 1)+".log")
				logfile, err := os.Create(logpath)
				if err != nil {
					continue
				}

				RedirectStderr(logfile)

				b.file.Close()
				b.file = logfile
				b.time = now

			}
			b.file.Write(b.format.Format(record))
			b.buffer += string(b.format.Format(record))
			if b.cnt == 100 {
				go b.handler(b.buffer)
				b.cnt = 0
				b.buffer = ""
			}

		case <-b.StopCh:
			break
		}
	}
}
