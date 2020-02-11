package log

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type BerithLogBatch struct {
	ch           chan *Record
	file         *os.File
	time         time.Time
	rotatePeriod time.Duration
	logdir       string
	format       Format
	StopCh       chan interface{}
}

func NewBerithLogBatch(ch chan *Record, logdir string, rotatePeriod time.Duration, format Format) *BerithLogBatch {
	return &BerithLogBatch{
		ch:           ch,
		logdir:       logdir,
		rotatePeriod: rotatePeriod,
		format:       format,
		StopCh:       make(chan interface{}),
	}
}

func (b *BerithLogBatch) Loop() {
	for {
		select {
		case record := <-b.ch:
			if b.file == nil || time.Now().Sub(b.time) >= b.rotatePeriod {

				if err := os.MkdirAll(b.logdir, 0700); err != nil {
					fmt.Println(err.Error())
					continue
				}
				now := time.Now()
				logpath := filepath.Join(b.logdir, strings.Replace(now.Format("060102150405.00"), ".", "", 1)+".log")
				logfile, err := os.Create(logpath)

				if err != nil {
					fmt.Println(err.Error())
					continue
				}

				b.file.Close()
				b.file = logfile
				b.time = now

			}
			b.file.Write(b.format.Format(record))
		case <-b.StopCh:
			break
		}
	}
}
