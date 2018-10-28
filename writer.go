// Copyright 2018 xgfone
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package miss

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
)

var fileFlag = os.O_CREATE | os.O_APPEND | os.O_WRONLY

// LevelWriter supports not only io.Writer but also WriteLevel.
type LevelWriter interface {
	io.Writer

	WriteLevel(level Level, bs []byte) (n int, err error)
}

// MayWriteLevel firstly tries to call the method WriteLevel to write the data.
// Or use io.Writer to write it.
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

type levelWriterFunc func(Level, []byte) (int, error)

func (l levelWriterFunc) Write(p []byte) (n int, err error) {
	return 0, errors.New("only support WriteLevel")
}

func (l levelWriterFunc) WriteLevel(lvl Level, p []byte) (n int, err error) {
	return l(lvl, p)
}

// LevelWriterFunc converts a function to LevelWriter.
func LevelWriterFunc(f func(Level, []byte) (int, error)) LevelWriter {
	return levelWriterFunc(f)
}

// LevelFilterWriter filters the logs whose level is less than lvl.
func LevelFilterWriter(lvl Level, w io.Writer) LevelWriter {
	return LevelWriterFunc(func(l Level, p []byte) (int, error) {
		if l < lvl {
			return 0, nil
		}
		return w.Write(p)
	})
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

	f, err := os.OpenFile(path, fileFlag, _mode)
	if err != nil {
		return nil, nil, err
	}
	return SafeWriter(f), f, nil
}

// ReopenWriter returns a writer that can be closed then re-opened,
// which is used for logrotate typically.
//
// Notice: it used SafeWriter to wrap the writer, so it's thread-safe.
func ReopenWriter(factory func() (w io.WriteCloser, reopen <-chan bool, err error)) (io.Writer, error) {
	w, reopen, err := factory()
	if err != nil {
		return nil, err
	}

	close := func() (int, error) {
		if w != nil {
			w.Close()
		}
		w = nil
		reopen = nil
		return 0, err
	}

	writer := WriterFunc(func(p []byte) (int, error) {
		if reopen == nil {
			if w, reopen, err = factory(); err != nil {
				return close()
			}
		}

		select {
		case <-reopen:
			w.Close()
			if w, reopen, err = factory(); err != nil {
				return close()
			}
		default:
		}
		return w.Write(p)
	})
	return SafeWriter(writer), nil
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

// SizedRotatingFileWriter returns a new file writer with rotating
// based on the size of the file.
//
// It is thread-safe for concurrent writes.
//
// The default permission of the log file is 0644.
func SizedRotatingFileWriter(filename string, size, count int,
	mode ...os.FileMode) (io.WriteCloser, error) {

	var _mode os.FileMode = 0644
	if len(mode) > 0 && mode[0] > 0 {
		_mode = mode[0]
	}

	w := sizedRotatingFile{
		filename:    filename,
		filePerm:    _mode,
		maxSize:     size,
		backupCount: count,
	}

	if err := w.open(); err != nil {
		return nil, err
	}
	return &w, nil
}

// sizedRotatingFile is a rotating logging handler based on the size.
type sizedRotatingFile struct {
	sync.Mutex
	file *os.File

	filePerm    os.FileMode
	filename    string
	maxSize     int
	backupCount int
	nbytes      int
}

func (f *sizedRotatingFile) Close() (err error) {
	f.Lock()
	err = f.close()
	f.Unlock()
	return
}

func (f *sizedRotatingFile) Write(data []byte) (n int, err error) {
	f.Lock()
	defer f.Unlock()

	if f.file == nil {
		return 0, errors.New("the file has been closed")
	}

	if f.nbytes+len(data) > f.maxSize {
		if err = f.doRollover(); err != nil {
			return
		}
	}

	if n, err = f.file.Write(data); err != nil {
		return
	}

	f.nbytes += n
	return
}

func (f *sizedRotatingFile) open() (err error) {
	file, err := os.OpenFile(f.filename, fileFlag, f.filePerm)
	if err != nil {
		return
	}

	info, err := file.Stat()
	if err != nil {
		return
	}

	f.nbytes = int(info.Size())
	f.file = file
	return
}

func (f *sizedRotatingFile) close() (err error) {
	err = f.file.Close()
	f.file = nil
	return
}

func (f *sizedRotatingFile) doRollover() (err error) {
	if f.backupCount > 0 {
		if err = f.close(); err != nil {
			return fmt.Errorf("Rotating: close failed: %s", err)
		}

		if !fileIsExist(f.filename) {
			return nil
		} else if n, err := fileSize(f.filename); err != nil {
			return fmt.Errorf("Rotating: failed to get the size: %s", err)
		} else if n == 0 {
			return nil
		}

		for _, i := range Range(f.backupCount-1, 0, -1) {
			sfn := fmt.Sprintf("%s.%d", f.filename, i)
			dfn := fmt.Sprintf("%s.%d", f.filename, i+1)
			if fileIsExist(sfn) {
				if fileIsExist(dfn) {
					os.Remove(dfn)
				}
				if err = os.Rename(sfn, dfn); err != nil {
					return fmt.Errorf("Rotating: failed to rename %s -> %s: %s",
						sfn, dfn, err)
				}
			}
		}
		dfn := f.filename + ".1"
		if fileIsExist(dfn) {
			if err = os.Remove(dfn); err != nil {
				return fmt.Errorf("Rotating: failed to remove %s: %s", dfn, err)
			}
		}
		if fileIsExist(f.filename) {
			if err = os.Rename(f.filename, dfn); err != nil {
				return fmt.Errorf("Rotating: failed to rename %s -> %s: %s",
					f.filename, dfn, err)
			}
		}
		err = f.open()
	}
	return
}

func fileIsExist(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// fileSize returns the size of the file as how many bytes.
func fileSize(fp string) (int64, error) {
	f, e := os.Stat(fp)
	if e != nil {
		return 0, e
	}
	return f.Size(), nil
}
