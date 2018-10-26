# miss

Package miss provides an flexible, extensible and powerful logging management tool based on the level, which has done the better balance between the flexibility and the performance.

miss is meaning:
1. **love**: Because of loving it, I miss it.
2. **flexible and extensible**: Something can be customized according to demand, so they are missing.
3. **no any third-party dependencies**: for the core package, you don't care any other packages, including the third-party.

It is inspired by [log15](https://github.com/inconshreveable/log15), [logrus](https://github.com/sirupsen/logrus), [go-kit](https://github.com/go-kit/kit).


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


## Example

1. Prepare a writer having implemented `io.Writer`, such as `os.Stdout`.
2. Create an encoder.
3. Create a logger with the encoder.
4. Output the log.

```go
package main

import (
	"os"

	"github.com/xgfone/miss"
)

func main() {
	conf := miss.EncoderConfig{IsLevel: true, IsTime: true}
	encoder := miss.KvTextEncoder(os.Stdout, conf)
	logger := miss.New(encoder).Level(miss.WARN)

	logger.Info("don't output")
	logger.Error("will output", "key", "value")
	// Output:
	// t=2018-10-25T10:46:22.0035694+08:00 lvl=ERROR key=value msg=will output
}
```

**Notice:**

`miss` is based on the level, and the log output interfaces is **`func(string, ...interface{}) error`**, the meaning of the arguments of which is decided by the encoder. See below.


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

The core package provides three kinds of encoders: the text encoder based on Key-Value `KvTextEncoder`, the text encoder based on Format `FmtTextEncoder` and the json encoder based on Key-Value `KvStdJSONEncoder` and `KvSimpleJSONEncoder`.

For the encoders based on Format, the arguments of the log output function, such as `Info()`, are the same as those of `fmt.Sprintf()`. For the encoders based on Key-Value, but, the first argument is the log description, and the rests are the key-value pairs, the number of which are even, for example, `logger.Info("log description", "key1", "value1", "key2", "value2")`.

```go
kvlog := miss.New(miss.KvTextEncoder(os.Stdout))
kvlog.Info("creating connection", "host", "127.0.0.1", "port", 80)

fmtlog := miss.New(miss.FmtTextEncoder(os.Stdout))
kvlog.Info("creating connection to %s:%d", "127.0.0.1", 80)
```

#### `LevelFilterEncoder` and `MultiEncoder`

You can use `LevelFilterEncoder` to filter some logs by the level, for example,

```go
encoders := ["kvtext", "kvjson"]
textenc := miss.KvTextEncoder(os.Stdout)
jsonenc := miss.KvSimpleJSONEncoder(os.Stderr)

textenc = miss.LevelFilterEncoder(miss.INFO, textenc)
jsonenc = miss.LevelFilterEncoder(miss.ERROR, jsonenc)

logger := miss.New(miss.MultiEncoder(textenc, jsonenc))

if err := logger.Info("only output to stdout"); err != nil {
    for i, e := range err.(miss.MultiError) {
        fmt.Printf("%s: %s\n", encoders[i], e.Error())
    }
}

if err := logger.Error("output to stdout & stderr); err != nil {
    for i, e := range err.(miss.MultiError) {
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
writer := miss.MultiWriter(os.Stdout, miss.FileWriter("/path/to/file"))
encoder := miss.KvTextEncoder(writer)
logger := miss.New(encoder)

logger.Info("output to stdout and file")
```


### Lazy evaluation

If the type of a certain value is `Valuer`, the default encoder engine will call it and encode the returned result. For example,
```go
logger := miss.New("hello", func() (interface{}, error) { return "world", nil })
```
or
```go
logger.Info("hello %v", func() (interface{}, error) { return "world", nil })
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
