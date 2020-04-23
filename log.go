package main

import (
	"os"
	"strconv"

	"github.com/go-sdk/logx"
)

type Log struct {
	i    int
	name string
}

func (log *Log) Printf(s string, v ...interface{}) {
	logx.Debugf("["+log.name+":"+strconv.FormatInt(int64(log.i), 10)+"] "+s, v...)
}

// --------------------------------------------------------------------------------

func exit(s string, e interface{}, v ...interface{}) {
	if d {
		l := logx.Caller(6)
		if e != nil {
			l = l.WithField("error", e)
		}
		l.Errorf(s, v...)
	} else {
		logx.Errorf(s, v...)
	}
	os.Exit(1)
}
