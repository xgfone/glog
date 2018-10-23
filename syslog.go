// +build !windows,!plan9

// This file refers to github.com/inconshreveable/log15#syslog.go

package miss

import (
	"bytes"
	"errors"
	"io"
	"log/syslog"
)

type syslogWriter struct {
	w *syslog.Writer
}

func (s syslogWriter) Write(p []byte) (int, error) {
	return 0, errors.New("only support WriteLevel")
}

func (s syslogWriter) WriteLevel(level Level, p []byte) (n int, err error) {
	v := string(bytes.TrimSpace(p))
	switch level {
	case FATAL:
		err = s.w.Emerg(v)
	case PANIC:
		err = s.w.Crit(v)
	case ERROR:
		err = s.w.Err(v)
	case WARN:
		err = s.w.Warning(v)
	case INFO:
		err = s.w.Info(v)
	default:
		err = s.w.Debug(v)
	}

	if err == nil {
		n = len(p)
	}
	return
}

// SyslogWriter opens a connection to the system syslog daemon
// by calling syslog.New and writes all logs to it.
func SyslogWriter(priority syslog.Priority, tag string) (io.Writer,
	io.Closer, error) {
	w, err := syslog.New(priority, tag)
	if err != nil {
		return nil, nil, err
	}
	return syslogWriter{w}, w, nil
}

// SyslogNetWriter opens a connection to a log daemon over the network
// and writes all logs to it.
func SyslogNetWriter(net, addr string, priority syslog.Priority,
	tag string) (io.Writer, io.Closer, error) {
	w, err := syslog.Dial(net, addr, priority, tag)
	if err != nil {
		return nil, nil, err
	}
	return syslogWriter{w}, w, nil
}

func (m muster) SyslogWriter(priority syslog.Priority,
	tag string) (io.Writer, io.Closer) {
	return must(SyslogWriter(priority, tag))
}

func (m muster) SyslogNetWriter(net, addr string, priority syslog.Priority,
	tag string) (io.Writer, io.Closer) {
	return must(SyslogNetWriter(net, addr, priority, tag))
}
