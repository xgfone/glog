package miss

const (
	traceNameS = "TRACE"
	debugNameS = "DEBUG"
	infoNameS  = "INFO"
	warnNameS  = "WARN"
	errorNameS = "ERROR"
	panicNameS = "PANIC"
	fatalNameS = "FATAL"

	unknownNameS = "UNKNOWN"
)

var (
	traceNameB = []byte(traceNameS)
	debugNameB = []byte(debugNameS)
	infoNameB  = []byte(infoNameS)
	warnNameB  = []byte(warnNameS)
	errorNameB = []byte(errorNameS)
	panicNameB = []byte(panicNameS)
	fatalNameB = []byte(fatalNameS)

	unknownNameB = []byte(unknownNameS)
)

// Predefine some levels
const (
	TRACE Level = iota // It will output the log unconditionally.
	DEBUG
	INFO
	WARN
	ERROR
	PANIC
	FATAL
)

// Level represents a level.
type Level int

// String returns the string representation.
func (l Level) String() string {
	switch l {
	case TRACE:
		return traceNameS
	case DEBUG:
		return debugNameS
	case INFO:
		return infoNameS
	case WARN:
		return warnNameS
	case ERROR:
		return errorNameS
	case PANIC:
		return panicNameS
	case FATAL:
		return fatalNameS
	default:
		return unknownNameS
	}
}

// Bytes returns the []byte representation.
func (l Level) Bytes() []byte {
	switch l {
	case TRACE:
		return traceNameB
	case DEBUG:
		return debugNameB
	case INFO:
		return infoNameB
	case WARN:
		return warnNameB
	case ERROR:
		return errorNameB
	case PANIC:
		return panicNameB
	case FATAL:
		return fatalNameB
	default:
		return unknownNameB
	}
}
