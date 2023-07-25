package main

import (
	"strings"

	"github.com/caticat/go_game_server/plog"
)

type PLogWriter struct {
	m_cha chan string
}

func NewPLogWriter() *PLogWriter {
	t := &PLogWriter{
		m_cha: make(chan string, PLOG_CHAN_LEN),
	}

	go t.run()

	return t
}

func (t *PLogWriter) Write(b []byte) (int, error) {
	t.m_cha <- string(b)

	return len(b), nil
}

func (t *PLogWriter) run() {
	logData := getLogData()
	logLast := getLogLast()

	for l := range t.m_cha {
		data, err := logData.Get()
		if err != nil {
			plog.Error(err)
			return
		}

		data += l
		if len(data) >= PLOG_MAX_SIZE {
			r := []rune(data)
			r = r[len(r)/2:] // 日志量直接减半
			data = string(r)
		}
		if err = logData.Set(data); err != nil {
			plog.Error(err)
			return
		}

		// 最后一行日志单独显示在主界面中
		if pos := strings.LastIndex(data[:(len(data)-1)], "\n"); pos >= 0 {
			logLast.Set(data[pos:(len(data) - 1)])
		}
	}
}
