# glog

Package glog provides an flexible, extensible and powerful logging management tool based on the level, which has done the better balance between the flexibility and the performance. It is inspired by [log15](https://github.com/inconshreveable/log15), [logrus](https://github.com/sirupsen/logrus), [go-kit](https://github.com/go-kit/kit).

See the [GoDoc](https://godoc.org/github.com/xgfone/glog).


## Prerequisite

Now `glog` requires Go `1.9+`.


## Basic Principle

- The better performance
- Flexible, extensible, and powerful
- No any third-party dependencies


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
type Logger interface {
	// Depth returns a new Logger with the stack depth.
	//
	// stackDepth is the calling depth of the logger, which will be passed to
	// the encoder. The default depth is the global variable DefaultLoggerDepth
	// for the new Logger.
	//
	// It should be used typically when you wrap the logger. For example,
	//
	//   logger := glog.New(glog.KvTextEncoder(os.Stdout))
	//   logger = logger.Depth(logger.GetDepth() + 1)
	//
	//   func Debug(m string, args ...interface{}) { logger.Debug(m, args...) }
	//   func Info(m string, args ...interface{}) { logger.Debug(m, args...) }
	//   func Warn(m string, args ...interface{}) { logger.Debug(m, args...) }
	//   ...
	//
	Depth(stackDepth int) Logger

	// Level returns a new Logger with the new level.
	Level(level Level) Logger

	// Encoder returns a new logger with the new encoder.
	Encoder(encoder Encoder) Logger

	// Ctx returns a new logger with the new contexts.
	Cxt(ctxs ...interface{}) Logger

	// Writer returns the underlying writer, which is the convenient function of
	// GetEncoder().Writer().
	Writer() io.Writer
	GetDepth() int
	GetLevel() Level
	GetEncoder() Encoder

	Trace(msg string, args ...interface{}) error
	Debug(msg string, args ...interface{}) error
	Info(msg string, args ...interface{}) error
	Warn(msg string, args ...interface{}) error
	Error(msg string, args ...interface{}) error
	Panic(msg string, args ...interface{}) error
	Fatal(msg string, args ...interface{}) error
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

	"github.com/xgfone/glog"
)

func main() {
	conf := glog.EncoderConfig{IsLevel: true, IsTime: true}
	encoder := glog.KvTextEncoder(os.Stdout, conf)
	logger := glog.New(encoder).Level(glog.WARN)

	logger.Info("don't output")
	logger.Error("will output", "key", "value")
	// Output:
	// t=2018-10-25T10:46:22.0035694+08:00 lvl=ERROR key=value msg=will output
}
```

Or you can use the convenient function `SimpleLogger(level, log_file_path string)`. If `log_file_path` is `""`, it will use `os.Stdout` as the output writer.

```go
package main

import (
	"os"

	"github.com/xgfone/glog"
)

func main() {
	logger, _, _ := glog.SimpleLogger("info", "")

	logger.Info("don't output")
	logger.Error("will output", "key", "value")
	// Output:
	// t=2018-10-25T10:46:22.0035694+08:00 lvl=ERROR key=value msg=will output
}
```

**Notice:**

`glog` is based on the level, and the log output interfaces is **`func(string, ...interface{}) error`**, the meaning of the arguments of which is decided by the encoder. See below.

Furthermore, `glog` has built in a global logger, which is equal to `glog.New(glog.FmtTextEncoder(os.Stdout, glog.EncoderConfig{IsLevel: true, IsTime: true}))`, and you can use the functions as follow:
```go
SetGlobalLogger(newLogger Logger)
GetGlobalLogger() Logger

WithLevel(level Level) Logger
WithEncoder(encoder Encoder) Logger
WithCtx(ctxs ...interface{}) Logger
WithDepth(depth int) Logger

GetDepth() int
GetLevel() Level
GetWriter() Writer // It's the short for GetEncode().Writer().
GetEncoder() Encoder

Trace(msg string, args ...interface{}) error
Debug(msg string, args ...interface{}) error
Info(msg string, args ...interface{}) error
Warn(msg string, args ...interface{}) error
Error(msg string, args ...interface{}) error
Panic(msg string, args ...interface{}) error
Fatal(msg string, args ...interface{}) error
```

If you prefer the logger without the error, you maybe use `LoggerWithoutError` converted by `ToLoggerWithoutError(Logger)` from `Logger` as follow:
```go
type LoggerWithoutError interface {
	Depth(stackDepth int) LoggerWithoutError
	Level(level Level) LoggerWithoutError
	Encoder(encoder Encoder) LoggerWithoutError
	Cxt(ctxs ...interface{}) LoggerWithoutError

	// Writer returns the underlying writer, which is the convenient function of
	// GetEncoder().Writer().
	Writer() io.Writer
	GetDepth() int
	GetLevel() Level
	GetEncoder() Encoder

	Trace(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Panic(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
}
```

### Inherit the context of the parent logger

```go
encoder := FmtTextEncoder(os.Stdout)
parent := logger.New(encoder).Ctx("parent")
child := parent.Ctx("child")
child.Info("hello %s", "world")
// Output:
// [parent][child] :=>: hello world
```

OR

```go
parent := logger.New("key1", "value1")
child := parent.New("key2", "value2").Encoder(KvTextEncoder(os.Stdout))
child.Info("hello world", "key3", "value3")
// Output:
// key1=value1 key2=value2 key3=value3 msg=hello world
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
	Encode(depth int, level Level, msg string, args []interface{}, ctx []interface{}) error
}
```

The core package provides three kinds of the implementations of the encoder: the text encoder based on Key-Value `KvTextEncoder`, the text encoder based on Format `FmtTextEncoder` and the json encoder based on Key-Value `KvStdJSONEncoder` and `KvSimpleJSONEncoder`.

For the encoders based on Format, the arguments of the log output function, such as `Info()`, are the same as those of `fmt.Sprintf()`. For the encoders based on Key-Value, but, the first argument is the log description, and the rests are the key-value pairs, the number of which are even, for example, `logger.Info("log description", "key1", "value1", "key2", "value2")`.

```go
kvlog := glog.New(glog.KvTextEncoder(os.Stdout))
kvlog.Info("creating connection", "host", "127.0.0.1", "port", 80)

fmtlog := glog.New(glog.FmtTextEncoder(os.Stdout))
kvlog.Info("creating connection to %s:%d", "127.0.0.1", 80)
```

#### `LevelFilterEncoder` and `MultiEncoder`

You can use `LevelFilterEncoder` to filter some logs by the level, for example,

```go
encoders := ["kvtext", "kvjson"]
textenc := glog.KvTextEncoder(os.Stdout)
jsonenc := glog.KvSimpleJSONEncoder(os.Stderr)

textenc = glog.LevelFilterEncoder(glog.INFO, textenc)
jsonenc = glog.LevelFilterEncoder(glog.ERROR, jsonenc)

logger := glog.New(glog.MultiEncoder(textenc, jsonenc))

if err := logger.Info("only output to stdout"); err != nil {
    for i, e := range err.(glog.MultiError) {
        fmt.Printf("%s: %s\n", encoders[i], e.Error())
    }
}

if err := logger.Error("output to stdout & stderr); err != nil {
    for i, e := range err.(glog.MultiError) {
        fmt.Printf("%s: %s\n", encoders[i], e.Error())
    }
}
```


### Writer

All implementing the interface `io.Writer` are a Writer.

There are some the built-in writers in the core package, such as `DiscardWriter`, `NetWriter`, `FileWriter`, `MultiWriter`, `FailoverWriter`, `SafeWriter`, `ChannelWriter`, `BufferedWriter`, `LevelFilterWriter`, `SyslogWriter`, `SyslogNetWriter`, `SizedRotatingFileWriter` and `Must`.


#### MultiWriter

For an encoder, you can output the result to more than one destination by using `MultiWriter`. For example, output the log to STDOUT and the file:

```go
writer := glog.MultiWriter(os.Stdout, glog.FileWriter("/path/to/file"))
encoder := glog.KvTextEncoder(writer)
logger := glog.New(encoder)

logger.Info("output to stdout and file")
```


### Lazy evaluation

If the type of a certain value is `Valuer`, the default encoder engine will call it and encode the returned result. For example,
```go
logger := glog.New("hello", func(d int, l Level) (interface{}, error) { return "world", nil })
```
or
```go
logger.Info("hello %v", func(d int, l Level) (interface{}, error) { return "world", nil })
```


## Performance

The log framework itself has no any performance costs.

There may be some performance costs below:
1. Use format arguments or Key-Value pairs when firing a log. For example, `logger.Info("hello %s", "world")` will allocate the 16-byte memory once for the encoder `FmtTextEncoder` , `logger.Info("hello world", "key", "value")` will allocate the 32-byte memory once for the encoder `KvTextEncoder`.
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
