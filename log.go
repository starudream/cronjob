package main

import (
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
