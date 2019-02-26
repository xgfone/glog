package logger

// Printfer is a printf log interface.
type Printfer interface {
	Printf(string, ...interface{})
}

// ToPrintfer converts Logger to Printfer.
//
// It will use logger.Info() to output the log by default. But you can change it
// by the level argument.
func ToPrintfer(logger Logger, level ...Level) Printfer {
	logger = logger.Depth(logger.GetDepth() + 2)
	logf := logger.Info
	if len(level) > 0 {
		switch l := level[0]; l {
		case LvlTrace:
			logf = logger.Trace
		case LvlDebug:
			logf = logger.Debug
		case LvlWarn:
			logf = logger.Warn
		case LvlError:
			logf = logger.Error
		case LvlPanic:
			logf = logger.Panic
		case LvlFatal:
			logf = logger.Fatal
		}
	}
	return printfer{logger: logf}
}

type printfer struct {
	logger func(string, ...interface{}) error
}

func (p printfer) Printf(format string, args ...interface{}) {
	p.logger(format, args...)
}
