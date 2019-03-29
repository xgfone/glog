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

package logger

import (
	"fmt"
	"strings"
)

const (
	traceNameS   = "TRACE"
	debugNameS   = "DEBUG"
	infoNameS    = "INFO"
	warnNameS    = "WARN"
	errorNameS   = "ERROR"
	panicNameS   = "PANIC"
	fatalNameS   = "FATAL"
	unknownNameS = "UNKNOWN"

	traceShortNameS   = "T"
	debugShortNameS   = "D"
	infoShortNameS    = "I"
	warnShortNameS    = "W"
	errorShortNameS   = "E"
	panicShortNameS   = "P"
	fatalShortNameS   = "F"
	unknownShortNameS = "U"
)

var (
	traceNameB   = []byte(traceNameS)
	debugNameB   = []byte(debugNameS)
	infoNameB    = []byte(infoNameS)
	warnNameB    = []byte(warnNameS)
	errorNameB   = []byte(errorNameS)
	panicNameB   = []byte(panicNameS)
	fatalNameB   = []byte(fatalNameS)
	unknownNameB = []byte(unknownNameS)

	traceShortNameB   = []byte(traceShortNameS)
	debugShortNameB   = []byte(debugShortNameS)
	infoShortNameB    = []byte(infoShortNameS)
	warnShortNameB    = []byte(warnShortNameS)
	errorShortNameB   = []byte(errorShortNameS)
	panicShortNameB   = []byte(panicShortNameS)
	fatalShortNameB   = []byte(fatalShortNameS)
	unknownShortNameB = []byte(unknownShortNameS)
)

// Predefine some levels
const (
	LvlTrace Level = iota // It will output the log unconditionally.
	LvlDebug
	LvlInfo
	LvlWarn
	LvlError
	LvlPanic
	LvlFatal
)

// Level represents a level.
type Level int32

// String returns the string representation.
func (l Level) String() string {
	switch l {
	case LvlTrace:
		return traceNameS
	case LvlDebug:
		return debugNameS
	case LvlInfo:
		return infoNameS
	case LvlWarn:
		return warnNameS
	case LvlError:
		return errorNameS
	case LvlPanic:
		return panicNameS
	case LvlFatal:
		return fatalNameS
	default:
		return unknownNameS
	}
}

// ShortString returns the short string representation.
func (l Level) ShortString() string {
	switch l {
	case LvlTrace:
		return traceShortNameS
	case LvlDebug:
		return debugShortNameS
	case LvlInfo:
		return infoShortNameS
	case LvlWarn:
		return warnShortNameS
	case LvlError:
		return errorShortNameS
	case LvlPanic:
		return panicShortNameS
	case LvlFatal:
		return fatalShortNameS
	default:
		return unknownShortNameS
	}
}

// Bytes returns the []byte representation.
func (l Level) Bytes() []byte {
	switch l {
	case LvlTrace:
		return traceNameB
	case LvlDebug:
		return debugNameB
	case LvlInfo:
		return infoNameB
	case LvlWarn:
		return warnNameB
	case LvlError:
		return errorNameB
	case LvlPanic:
		return panicNameB
	case LvlFatal:
		return fatalNameB
	default:
		return unknownNameB
	}
}

// ShortBytes returns the short []byte representation.
func (l Level) ShortBytes() []byte {
	switch l {
	case LvlTrace:
		return traceShortNameB
	case LvlDebug:
		return debugShortNameB
	case LvlInfo:
		return infoShortNameB
	case LvlWarn:
		return warnShortNameB
	case LvlError:
		return errorShortNameB
	case LvlPanic:
		return panicShortNameB
	case LvlFatal:
		return fatalShortNameB
	default:
		return unknownShortNameB
	}
}

// WriteTo writes the level into out.
func (l Level) WriteTo(out Writer, short bool) (n int, err error) {
	if short {
		return out.Write(l.ShortBytes())
	}
	return out.Write(l.Bytes())
}

// NameToLevel returns the Level by the name, which is case Insensitive.
//
// It supports the full or short name, but panic if the name is unknown.
//
// Notice: WARNING is the alias of WARN.
func NameToLevel(name string) Level {
	switch strings.ToUpper(name) {
	case traceNameS, traceShortNameS:
		return LvlTrace
	case debugNameS, debugShortNameS:
		return LvlDebug
	case infoNameS, infoShortNameS:
		return LvlInfo
	case warnNameS, warnShortNameS, "WARNING":
		return LvlWarn
	case errorNameS, errorShortNameS:
		return LvlError
	case panicNameS, panicShortNameS:
		return LvlPanic
	case fatalNameS, fatalShortNameS:
		return LvlFatal
	default:
		panic(fmt.Errorf("unknown level name '%s'", name))
	}
}
