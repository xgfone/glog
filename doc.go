// Package miss provides an flexible, extensible and powerful logging
// management tool based on the level, which has done the better balance
// between the flexibility and the performance.
//
// Meaning
//   - love: Because of loving it, I miss it.
//   - flexible and extensible: Something can be customized according to demand, so they are missing.
//   - no any third-party dependencies: for the core package, you don't care any other packages, including the third-party.
//
// Basic Principle
//   - The better performance
//   - Flexible, extensible, and powerful
//   - No any third-party dependencies
//
// Features
//   - A simple, easy-to-understand API
//   - No any third-party dependencies for the core package.
//   - A flexible and powerful interface supporting many encoders
//   - Child loggers which inherit and add their own private context
//   - Lazy evaluation of expensive operations
//   - Support any `io.Writer` and provied some advanced `io.Writer` implementations
//   - Built-in support for logging to files, syslog, and the network
//
// Encoder
//
// The core package provides three kinds of encoders:
//   - the text encoder based on Key-Value `KvTextEncoder`
//   - the text encoder based on Format `FmtTextEncoder`
//   - the json encoder based on Key-Value `KvStdJSONEncoder` and `KvSimpleJSONEncoder`
//
// For the encoders based on Format, the arguments of the log output function, such as `Info()`, are the same as those of `fmt.Sprintf()`. For the encoders based on Key-Value, howerer, the first argument is the log description, and the rests are the key-value pairs, the number of which are even, for example,
//  logger.Info("log description", "key1", "value1", "key2", "value2")
//
// You can use `LevelFilterEncoder` to filter some logs by the level, for example,
//
// Writer
//
// All implementing the interface `io.Writer` are a Writer.
//
// There are some the built-in writers in the core package, For example,
//
//   NetWriter    SyslogWriter    SizedRotatingFileWriter
//   FileWriter   ChannelWriter   LevelFilterWriter
//   SafeWriter   DiscardWriter   SyslogNetWriter
//   MultiWriter  BufferedWriter  FailoverWriter
//
// Performance
//
// The log framework itself has no any performance costs.
//
// There may be some performance costs below:
//   1. Use format arguments or Key-Value pairs when firing a log.
//      For example, `logger.Info("hello %s", "world")` will allocate
//      the 16-byte memory once for the encoder `FmtTextEncoder`,
//      `logger.Info("hello world", "key", "value")` will allocate
//      the 32-byte memory once for the encoder `KvTextEncoder`.
//   2. Encode the arguments to `io.Writer`. For `string` or `[]byte`,
//      there is no any performance cost, but for other types,
//      such as `int`, it maybe have once memory allocation.
//
package miss
