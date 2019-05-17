# logger [![Build Status](https://travis-ci.org/xgfone/logger.svg?branch=master)](https://travis-ci.org/xgfone/logger) [![GoDoc](https://godoc.org/github.com/xgfone/logger?status.svg)](http://godoc.org/github.com/xgfone/logger) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg?style=flat-square)](https://raw.githubusercontent.com/xgfone/logger/master/LICENSE)

Package `logger` provides an flexible, extensible and powerful logging management tool based on the level, which has done the better balance between the flexibility and the performance. It is inspired by [log15](https://github.com/inconshreveable/log15), [logrus](https://github.com/sirupsen/logrus), [go-kit](https://github.com/go-kit/kit).

See the [GoDoc](https://godoc.org/github.com/xgfone/logger).

**API has been stable.** The current is `v4`.


## Prerequisite

Now `logger` requires Go `1.7+`.


## Basic Principle

- The better performance
- Flexible, extensible, and powerful


## Features

- A simple, easy-to-understand API
- No any third-party dependencies for the core package.
- A flexible and powerful interface supporting many encoders, such as the `Key-Value` or `Format` style encoder
- Child loggers which inherit and add their own private context
- Lazy evaluation of expensive operations
- Support the native `io.Writer` as the output, and provied some advanced `io.Writer` implementations, such as `MultiWriter` and `LevelWriter`.
- Built-in support for logging to files, syslog, and the network


## `Logger`

```go
// LogGetter is an interface to return the inner information of Logger.
type LogGetter interface {
	GetName() string
	GetDepth() int
	GetLevel() Level
	GetEncoder() Encoder
}

// LogSetter is an interface to modify the inner information of Logger.
type LogSetter interface {
	SetName(name string)
	SetDepth(depth int)
	SetLevel(level Level) // It should be thread-safe.
	SetEncoder(encoder Encoder)
}

// LogWither is an interface to return a new Logger based on the current logger
// with the new argument.
type LogWither interface {
	WithName(name string) Logger
	WithLevel(level Level) Logger
	WithEncoder(encoder Encoder) Logger
	WithCxt(ctxs ...interface{}) Logger
	WithDepth(stackDepth int) Logger
}

// LogOutputter is an interface to emit the log.
type LogOutputter interface {
	Trace(msg string, args ...interface{}) error
	Debug(msg string, args ...interface{}) error
	Info(msg string, args ...interface{}) error
	Warn(msg string, args ...interface{}) error
	Error(msg string, args ...interface{}) error
	Panic(msg string, args ...interface{}) error
	Fatal(msg string, args ...interface{}) error
}

// Logger is an compositive logger interface.
type Logger interface {
	LogGetter
	LogSetter
	LogWither
	LogOutputter
}
```


## Example

1. Prepare a writer having implemented `io.Writer`, such as `os.Stdout`.
2. Create an encoder.
3. Create a logger with the encoder.
4. Output the log.

```go
package main

import (
	"os"

	"github.com/xgfone/logger"
)

func main() {
	encoder := logger.NewTextJSONEncoder(os.Stdout)
	log := logger.New(encoder).WithLevel(logger.LvlWarn)

	log.Info("don't output")
	log.Error("will output", "key", "value")

	// Output:
	// time=2019-05-17T15:34:00.9473464+08:00 level=ERROR key=value msg=will output
}
```

Or you can use the convenient function `SimpleLogger(level, log_file_path string)`. If `log_file_path` is `""`, it will use `os.Stdout` as the output writer.

```go
package main

import (
	"github.com/xgfone/logger"
)

func main() {
	log, _, _ := logger.SimpleLogger("warn", "")

	log.Info("don't output")
	log.Error("will output %s %s", "key", "value")

	// Output:
	// 2019-05-17T15:36:13.3431634+08:00 [ERROR]: will output key value
}
```

**Notice:**

`logger` is based on the level, and the log output interfaces is **`func(string, ...interface{}) error`**, the meaning of the arguments of which is decided by the encoder. See below.

Furthermore, `logger` has built in a global logger, which is equal to `logger.New(logger.NewFmtEncoder(os.Stdout))`, and you can use the functions as follow:
```go
SetGlobalLogger(newLogger Logger)
GetGlobalLogger() Logger

WithName(name string) Logger
WithLevel(level Level) Logger
WithEncoder(encoder Encoder) Logger
WithCtx(ctxs ...interface{}) Logger
WithDepth(depth int) Logger

GetName() string
GetDepth() int
GetLevel() Level
GetWriter() Writer // It's the short for GetEncode().Writer().
GetEncoder() Encoder

SetName(name string)
SetDepth(depth int)
SetLevel(level Level)
SetEncoder(encoder Encoder)

Trace(msg string, args ...interface{}) error
Debug(msg string, args ...interface{}) error
Info(msg string, args ...interface{}) error
Warn(msg string, args ...interface{}) error
Error(msg string, args ...interface{}) error
Panic(msg string, args ...interface{}) error
Fatal(msg string, args ...interface{}) error
```

**Suggestion:**
Use the global logger instead of the customized logger directly, such as `logger.Trace()`, `logger.Debug()`, `logger.Info()`, `logger.Warn()`, `logger.Error()`, `logger.Panic()`, `logger.Fatal()`. If you need to use a new logger, you can set the global logger to it by `logger.SetGlobalLogger()` on initializating the program.

If you prefer the logger without the error, you maybe use `NoErrorLogger` converted by `ToNoErrorLogger(Logger)` from `Logger` as follow:
```go
// LogOutputterWithoutError is the same as LogOutputter, but not return an error.
type LogOutputterWithoutError interface {
	Trace(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Panic(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
}

// NoErrorLogger is equal to Logger, but not returning the error.
type NoErrorLogger interface {
	LogGetter
	LogSetter
	LogOutputterWithoutError

	WithName(name string) NoErrorLogger
	WithDepth(stackDepth int) NoErrorLogger
	WithLevel(level Level) NoErrorLogger
	WithEncoder(encoder Encoder) NoErrorLogger
	WithCxt(ctxs ...interface{}) NoErrorLogger
}
```

### Inherit the context of the parent logger

```go
package main

import (
	"os"

	"github.com/xgfone/logger"
)

func main() {
	encoder := logger.NewFmtEncoder(os.Stdout)
	parent := logger.New(encoder).WithCxt("parent")
	child := parent.WithCxt("child")
	child.Info("hello %s", "world")

	// Output:
	// 2019-05-17T15:42:36.2185998+08:00 parent|child main.go:13 [INFO]: hello world
}
```


### Encoder

```go
type Encoder interface {
	// Reset the underlying writer.
	ResetWriter(Writer)

	// Return the underlying writer.
	//
	// Notice: only the most underlying encoder requires it. For the inner
	// encoder, such as FilterEncoder and MultiEncoder, it may be nil.
	// So, at the moment, the log information should be passed to the next encoder.
	Writer() Writer

	// Encode the log and write it into the underlying writer.
	Encode(Record) error
}
```

The core package provides three kinds of the implementations of the encoder: the text encoder based on Key-Value `NewTextJSONEncoder`, the text encoder based on Format `NewFmtEncoder`, and the json encoder based on Key-Value `NewStdJSONEncoder` and `NewSimpleJSONEncoder`.

For the encoders based on Format, the arguments of the log output function, such as `Info()`, are the same as those of `fmt.Sprintf()`. For the encoders based on Key-Value, but, the first argument is the log description, and the rests are the key-value pairs, the number of which are even, for example, `logger.Info("log description", "key1", "value1", "key2", "value2")`.

```go
kvlog := logger.New(logger.NewTextJSONEncoder(os.Stdout))
kvlog.Info("creating connection", "host", "127.0.0.1", "port", 80)

fmtlog := logger.New(logger.NewFmtEncoder(os.Stdout))
kvlog.Info("creating connection to %s:%d", "127.0.0.1", 80)
```

#### `LevelFilterEncoder` and `MultiEncoder`

You can use `LevelFilterEncoder` to filter some logs by the level, for example,

```go
encoders := ["kvtext", "kvjson"]
textenc := logger.NewTextJSONEncoder(os.Stdout)
jsonenc := logger.NewSimpleJSONEncoder(os.Stderr)

textenc = logger.LevelFilterEncoder(logger.LvlInfo, textenc)
jsonenc = logger.LevelFilterEncoder(logger.LvlError, jsonenc)

log := logger.New(logger.MultiEncoder(textenc, jsonenc))

if err := log.Info("only output to stdout"); err != nil {
    for i, e := range err.(logger.MultiError) {
        fmt.Printf("%s: %s\n", encoders[i], e.Error())
    }
}

if err := log.Error("output to stdout & stderr"); err != nil {
    for i, e := range err.(logger.MultiError) {
        fmt.Printf("%s: %s\n", encoders[i], e.Error())
    }
}
```


### Writer

All implementing the interface `io.Writer` are a Writer.

There are some the built-in writers in the core package, such as `DiscardWriter`, `NetWriter`, `FileWriter`, `MultiWriter`, `FailoverWriter`, `SafeWriter`, `ChannelWriter`, `BufferedWriter`, `LevelFilterWriter`, `SyslogWriter`, `SyslogNetWriter`, and `SizedRotatingFileWriter`.


#### MultiWriter

For an encoder, you can output the result to more than one destination by using `MultiWriter`. For example, output the log to STDOUT and the file:

```go
fileWriter, fileCloser, _ := logger.FileWriter("/path/to/file")
defer fileCloser.Close()

writer := logger.MultiWriter(os.Stdout, fileWriter)
encoder := logger.NewTextJSONEncoder(writer)
log := logger.New(encoder)

log.Info("output to stdout and file")
```


### Lazy evaluation

If the type of a certain value is `Valuer`, the default encoder engine will call it and encode the returned result. For example,
```go
lazy := func(r logger.Record) (interface{}, error) { return "world", nil }
log := logger.New(logger.NewTextJSONEncoder(os.Stdout)).WithCxt("hello", lazy)
```
or
```go
log.Info("hello %s", func(r logger.Record) (interface{}, error) { return "world", nil })
```


### Logger Cache
```go
// Initializing
SetEncoder(newEncoder)

// Usage
webLogger := GetLogger("web")
cliLogger := GetLogger("cli")

webLogger.Info("This is a web log")
cliLogger.Info("This is a cli log")
```

`GetLogger` will return the same logger for the same name, no matter how many times you call it.


## Performance

The log framework itself has no any performance costs.

There may be some performance costs below:
1. Use format arguments or Key-Value pairs when firing a log. For example, `logger.Info("hello %s", "world")` will allocate the 16-byte memory once for the encoder `FmtTextEncoder` , `logger.Info("hello world", "key", "value")` will allocate the 32-byte memory once for the encoder `KvTextEncoder`. For `NewFmtEncoder` and `NewTextJSONEncoder`, however, it will consume a little more memory because of `interface{}`.
2. Encode the arguments to `io.Writer`. For `string` or `[]byte`, there is no any performance cost, but for other types, such as `int`, it maybe have once memory allocation.


### Performance Test

The test program is from [go-loggers-bench](https://github.com/imkira/go-loggers-bench).

```
MacBook Pro(Retina, 13-inch, Mid 2014)
2.6 GHz Intel Core i5
8 GB 1600 MHz DDR3
macOS Mojave
```

#### TextNegative
|  test   | ops | ns/op | bytes/op | allocs/op
|---------|-----|-------|----------|-----------
| **BenchmarkMissTextNegative-4**      | **1000000000**  | **7.25 ns/op** | **0 B/op**   | **0 allocs/op**
| BenchmarkGokitTextNegative-4     | 300000000   | 26.8 ns/op | 32 B/op  | 1 allocs/op
| BenchmarkLog15TextNegative-4     | 10000000    | 723 ns/op  | 368 B/op | 3 allocs/op
| BenchmarkLogrusTextNegative-4    | 10000000000 | 1.93 ns/op | 0 B/op   | 0 allocs/op
| BenchmarkSeelogTextNegative-4    | 200000000   | 47.4 ns/op | 48 B/op  | 2 allocs/op
| BenchmarkZerologTextNegative-4   | 2000000000  | 4.39 ns/op | 0 B/op   | 0 allocs/op
| BenchmarkGologgingTextNegative-4 | 100000000   | 92.1 ns/op | 144 B/op | 2 allocs/op


#### TextPositive
|  test   | ops | ns/op | bytes/op | allocs/op
|---------|-----|-------|----------|-----------
| **BenchmarkMissTextPositive-4**      | **20000000** | **372 ns/op**  | **48 B/op**  | **1 allocs/op**
| BenchmarkLog15TextPositive-4     | 1000000  | 5125 ns/op | 856 B/op | 14 allocs/op
| BenchmarkGokitTextPositive-4     | 10000000 | 846 ns/op  | 256 B/op | 4 allocs/op
| BenchmarkSeelogTextPositive-4    | 2000000  | 3313 ns/op | 440 B/op | 11 allocs/op
| BenchmarkLogrusTextPositive-4    | 2000000  | 4433 ns/op | 448 B/op | 12 allocs/op
| BenchmarkZerologTextPositive-4   | 30000000 | 288 ns/op  | 0 B/op   | 0 allocs/op
| BenchmarkGologgingTextPositive-4 | 10000000 | 1093 ns/op | 920 B/op | 15 allocs/op


#### JSONNegative
|  test   | ops | ns/op | bytes/op | allocs/op
|---------|-----|-------|----------|-----------
| **BenchmarkMissJSONNegative-4**    | **200000000**  | **41.9 ns/op** | **96 B/op**  | **1 allocs/op**
| BenchmarkLog15JSONNegative-4   | 10000000   | 813 ns/op  | 560 B/op | 5 allocs/op
| BenchmarkGokitJSONNegative-4   | 200000000  | 45.0 ns/op | 128 B/op | 1 allocs/op
| BenchmarkLogrusJSONNegative-4  | 20000000   | 367 ns/op  | 480 B/op | 4 allocs/op
| BenchmarkZerologJSONNegative-4 | 1000000000 | 8.80 ns/op | 0 B/op   | 0 allocs/op


#### JSONPositive
|  test   | ops | ns/op | bytes/op | allocs/op
|---------|-----|-------|----------|-----------
| **BenchmarkMissJSONPositive-4**    | **10000000** | **1181 ns/op**  | **640 B/op**  | **10 allocs/op**
| BenchmarkLog15JSONPositive-4   | 500000   | 11108 ns/op | 2256 B/op | 32 allocs/op
| BenchmarkGokitJSONPositive-4   | 3000000  | 2726 ns/op  | 1552 B/op | 24 allocs/op
| BenchmarkLogrusJSONPositive-4  | 1000000  | 9273 ns/op  | 1843 B/op | 30 allocs/op
| BenchmarkZerologJSONPositive-4 | 20000000 | 397 ns/op   | 0 B/op    | 0 allocs/op

**NOTICE:** For pursuing the extreme performance, you maybe see [zerolog](https://github.com/rs/zerolog).
