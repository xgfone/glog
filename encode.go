package miss

import (
	"fmt"
	"io"
	"time"
)

// The separators of the KV and the KV pair.
const (
	TextKVSep     = "="
	TextKVPairSep = " "
)

// Some key names. You can modify them to redefine them.
const (
	LevelKey = "lvl"
	TimeKey  = "t"
	MsgKey   = "msg"
)

// ErrKeyValueNum will be used when the number of key-values is not even.
var ErrKeyValueNum = fmt.Errorf("the number of key-values must be even")

// Encoder is a log encoder.
//
// Notice: if the encoder implementation supports the level when writing data,
// it should firstly decide whether the writer is LevelWriter and use WriteLevel
// to write the log, not Write.
type Encoder interface {
	Encode(level Level, msg string, args []interface{}, ctx []interface{}) error
}

type encoderFunc func(Level, string, []interface{}, []interface{}) error

func (e encoderFunc) Encode(l Level, m string, args, ctx []interface{}) error {
	return e(l, m, args, ctx)
}

// EncoderFunc converts a function to an Encoder.
func EncoderFunc(f func(Level, string, []interface{}, []interface{}) error) Encoder {
	return encoderFunc(f)
}

// FilterEncoder returns an encoder that only forwards logs
// to the wrapped encoder if the given function evaluates true.
//
// For example, filter those logs that the level is less than ERROR.
//
//    FilterEncoder(func(lvl Level, msg string, args []interface{},
//                       ctxs []interface{}) bool {
//        return level >= ERROR
//    })
//
func FilterEncoder(f func(Level, string, []interface{}, []interface{}) bool,
	encoder Encoder) Encoder {
	return EncoderFunc(func(l Level, m string, args []interface{},
		ctxs []interface{}) error {
		if f(l, m, args, ctxs) {
			return encoder.Encode(l, m, args, ctxs)
		}
		return nil
	})
}

// LevelFilterEncoder returns an encoder that only writes records which are
// greater than the given verbosity level to the wrapped Handler.
//
// For example, to only output Error/PANIC/FATAL logs:
//
//     miss.LevelFilterEncoder(miss.ERROR, miss.KvTextEncoder(os.Stdout))
//
func LevelFilterEncoder(level Level, encoder Encoder) Encoder {
	return FilterEncoder(func(l Level, m string, args, ctxs []interface{}) bool {
		return l >= level
	}, encoder)
}

// NothingEncoder returns an encoder that does nothing.
func NothingEncoder() Encoder {
	return EncoderFunc(func(l Level, m string, args, ctx []interface{}) error {
		return nil
	})
}

// EncoderConfig configures the encoder.
type EncoderConfig struct {
	Slice []interface{}
	Map   map[string]interface{}

	// If true, the encoder disable appending a newline.
	NotNewLine bool

	// TimeLayout is used to format time.Time.
	//
	// The default is time.RFC3339Nano.
	TimeLayout string

	// If ture, the encoder will encode the current time.
	IsTime bool

	// If ture, the encoder will encode the level.
	IsLevel bool

	// For the Key-Value encoder, it represents the key name of the time.
	// The global constant, TimeKey, will be used by default.
	TimeKey string

	// For the Key-Value encoder, it represents the key name of the level.
	// The global constant, LvlKey, will be used by default.
	LevelKey string

	// For the Key-Value encoder, it represents the key name of the message.
	// The global constant, MsgKey, will be used by default.
	MsgKey string

	// The separator between key and value, such as "=".
	// The global constant, TextKVSep, will be used by default.
	TextKVSep string

	// The separator between the key-value pairs, such as " ".
	// The global constant, TextKVPairSep, will be used by default.
	TextKVPairSep string
}

func (ec EncoderConfig) init() EncoderConfig {
	if ec.TimeLayout == "" {
		ec.TimeLayout = time.RFC3339Nano
	}

	if ec.Slice == nil {
		ec.Slice = make([]interface{}, 0)
	}

	if ec.Map == nil {
		ec.Map = make(map[string]interface{})
	}

	return ec
}

// KvTextEncoder returns a text encoder based on the key-value pair,
// which will output the result into out.
//
// Notice: This encoder supports LevelWriter.
func KvTextEncoder(out io.Writer, conf ...EncoderConfig) Encoder {
	var c EncoderConfig

	if len(conf) > 0 {
		c = conf[0]

		if c.TimeKey == "" {
			c.TimeKey = TimeKey
		}
		if c.LevelKey == "" {
			c.LevelKey = LevelKey
		}
		if c.MsgKey == "" {
			c.MsgKey = MsgKey
		}

		if c.TextKVSep == "" {
			c.TextKVSep = TextKVSep
		}
		if c.TextKVPairSep == "" {
			c.TextKVPairSep = TextKVPairSep
		}

		if len(c.Slice) > 0 {
			c.Slice = append([]interface{}{}, c.Slice...)
		}

		if len(c.Map) > 0 {
			maps := make(map[string]interface{}, len(c.Map))
			for k, v := range c.Map {
				maps[k] = v
			}
			c.Map = maps
		}
	}

	return kvTextEncoder{out: out, conf: c.init()}
}

type kvTextEncoder struct {
	out  io.Writer
	conf EncoderConfig
}

func (t kvTextEncoder) Encode(l Level, m string, args, ctxs []interface{}) error {
	arglen := len(args)
	ctxlen := len(ctxs)
	if arglen%2 != 0 || ctxlen%2 != 0 {
		return ErrKeyValueNum
	}

	var err error
	var sep bool
	bs := bytesPools.Get()
	defer bytesPools.Put(bs)

	if t.conf.IsTime {
		var _bs [64]byte
		bs = append(bs, 't')
		bs = append(bs, t.conf.TextKVSep...)
		bs = append(bs, time.Now().AppendFormat(_bs[:0], t.conf.TimeLayout)...)
		sep = true
	}

	if t.conf.IsLevel {
		if sep {
			bs = append(bs, t.conf.TextKVPairSep...)
		}

		bs = append(bs, t.conf.LevelKey...)
		bs = append(bs, t.conf.TextKVSep...)
		bs = append(bs, l.Bytes()...)
		sep = true
	}

	for i := 0; i < ctxlen; i += 2 {
		if sep {
			bs = append(bs, t.conf.TextKVPairSep...)
		}

		if bs, err = WriteIntoBytes(bs, ctxs[i]); err != nil {
			return err
		}
		bs = append(bs, t.conf.TextKVSep...)
		if bs, err = WriteIntoBytes(bs, ctxs[i+1]); err != nil {
			return err
		}

		sep = true
	}

	for i := 0; i < arglen; i += 2 {
		if sep {
			bs = append(bs, t.conf.TextKVPairSep...)
		}

		if bs, err = WriteIntoBytes(bs, args[i]); err != nil {
			return err
		}
		bs = append(bs, t.conf.TextKVSep...)
		if bs, err = WriteIntoBytes(bs, args[i+1]); err != nil {
			return err
		}

		sep = true
	}

	if sep {
		bs = append(bs, t.conf.TextKVPairSep...)
	}

	bs = append(bs, t.conf.MsgKey...)
	bs = append(bs, t.conf.TextKVSep...)
	bs = append(bs, m...)

	if !t.conf.NotNewLine {
		bs = append(bs, '\n')
	}

	_, err = MayWriteLevel(t.out, l, bs)
	return err
}

// FmtTextEncoder returns a text encoder based on the % formatter,
// which will output the result into out.
//
// Notice: This encoder supports LevelWriter.
func FmtTextEncoder(out io.Writer, conf ...EncoderConfig) Encoder {
	var c EncoderConfig
	if len(conf) > 0 {
		c = conf[0]

		if len(c.Slice) > 0 {
			c.Slice = append([]interface{}{}, c.Slice...)
		}

		if len(c.Map) > 0 {
			maps := make(map[string]interface{}, len(c.Map))
			for k, v := range c.Map {
				maps[k] = v
			}
			c.Map = maps
		}
	}

	return fmtTextEncoder{out: out, conf: c.init()}
}

type fmtTextEncoder struct {
	out  io.Writer
	conf EncoderConfig
}

func (f fmtTextEncoder) Encode(l Level, m string, args, ctxs []interface{}) error {
	var err error
	var sep bool
	bs := bytesPools.Get()
	defer bytesPools.Put(bs)

	if f.conf.IsTime {
		var _bs [64]byte
		bs = append(bs, time.Now().AppendFormat(_bs[:0], f.conf.TimeLayout)...)
		sep = true
	}

	if f.conf.IsLevel {
		if sep {
			bs = append(bs, ' ')
		}
		bs = append(bs, l.Bytes()...)
		sep = true
	}

	ctxlen := len(ctxs)
	if ctxlen > 0 {
		if sep {
			bs = append(bs, ' ')
		}

		for _, v := range ctxs {
			bs = append(bs, '[')
			if bs, err = WriteIntoBytes(bs, v); err != nil {
				return err
			}
			bs = append(bs, ']')
		}

		sep = true
	}

	if sep {
		bs = append(bs, " :=>: "...)
	}

	bs = append(bs, fmt.Sprintf(m, args...)...)

	if !f.conf.NotNewLine {
		bs = append(bs, '\n')
	}

	_, err = MayWriteLevel(f.out, l, bs)
	return err
}