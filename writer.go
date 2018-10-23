package miss

import (
	"io"
	"net"
	"os"
	"sync"
)

// LevelWriter supports not only io.Writer but also WriteLevel.
type LevelWriter interface {
	io.Writer

	WriteLevel(level Level, bs []byte) (n int, err error)
}

// MayWriteLevel try
func MayWriteLevel(w io.Writer, level Level, bs []byte) (int, error) {
	if lw, ok := w.(LevelWriter); ok {
		return lw.WriteLevel(level, bs)
	}
	return w.Write(bs)
}

type writerFunc func([]byte) (int, error)

func (w writerFunc) Write(p []byte) (int, error) {
	return w(p)
}

// WriterFunc converts a function to io.Writer.
func WriterFunc(f func([]byte) (int, error)) io.Writer {
	return writerFunc(f)
}

// DiscardWriter returns a writer which will discard all input.
func DiscardWriter() io.Writer {
	return WriterFunc(func(p []byte) (int, error) {
		return len(p), nil
	})
}

// NetWriter opens a socket to the given address and writes the log
// over the connection.
func NetWriter(network, addr string) (io.Writer, io.Closer, error) {
	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, nil, err
	}

	return SafeWriter(conn), conn, nil
}

// FileWriter returns a writer which writes log records to the give file.
//
// If the path already exists, FileHook will append to the given file.
// If it does not, FileHook will create the file with mode 0644,
// but you can pass the second argument, mode, to modify it.
func FileWriter(path string, mode ...os.FileMode) (io.Writer, io.Closer, error) {
	var _mode os.FileMode = 0644
	if len(mode) > 0 {
		_mode = mode[0]
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, _mode)
	if err != nil {
		return nil, nil, err
	}
	return SafeWriter(f), f, nil
}

// MultiWriter writes one data to more than one destination.
func MultiWriter(outs ...io.Writer) io.Writer {
	return WriterFunc(func(p []byte) (n int, err error) {
		for _, out := range outs {
			if m, e := out.Write(p); e != nil {
				n = m
				err = e
			}
		}
		return
	})
}

// FailoverWriter writes all log records to the first handler specified,
// but will failover and write to the second handler if the first handler
// has failed, and so on for all handlers specified.
//
// For example, you might want to log to a network socket,
// but failover to writing to a file if the network fails,
// and then to standard out if the file write fails.
func FailoverWriter(outs ...io.Writer) io.Writer {
	return WriterFunc(func(p []byte) (n int, err error) {
		for _, out := range outs {
			if n, err = out.Write(p); err == nil {
				return
			}
		}
		return
	})
}

// SafeWriter is guaranteed that only a single writing operation
// can proceed at a time.
//
// It's necessary for thread-safe concurrent writes.
func SafeWriter(w io.Writer) io.Writer {
	var mu sync.Mutex
	return WriterFunc(func(p []byte) (int, error) {
		mu.Lock()
		defer mu.Unlock()
		return w.Write(p)
	})
}

// ChannelWriter writes all logs to the given channel.
//
// It blocks if the channel is full. Useful for async processing
// of log messages, it's used by BufferedWriter.
func ChannelWriter(ch chan<- []byte) io.Writer {
	return WriterFunc(func(p []byte) (int, error) {
		ch <- p
		return len(p), nil
	})
}

// BufferedWriter writes all records to a buffered channel of the given size
// which flushes into the wrapped handler whenever it is available for writing.
//
// Since these writes happen asynchronously, all writes to a BufferedWriter
// never return an error and any errors from the wrapped writer are ignored.
func BufferedWriter(bufSize int, w io.Writer) io.Writer {
	ch := make(chan []byte, bufSize)
	go func() {
		for bs := range ch {
			w.Write(bs)
		}
	}()
	return ChannelWriter(ch)
}

// Must object provides the following writer creation functions
// which instead of returning an error parameter only return a writer
// and panic on failure: FileWriter, NetWriter, SyslogWriter, SyslogNetWriter.
var Must muster

func must(w io.Writer, c io.Closer, err error) (io.Writer, io.Closer) {
	if err != nil {
		panic(err)
	}
	return w, c
}

type muster struct{}

func (m muster) FileWriter(path string, mode ...os.FileMode) (io.Writer, io.Closer) {
	return must(FileWriter(path, mode...))
}

func (m muster) NetWriter(network, addr string) (io.Writer, io.Closer) {
	return must(NetWriter(network, addr))
}
