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

// NameToLevel returns the Level by the name, which is case Insensitive.
//
// If the name is unknown, it will panic.
//
// Notice: WARNING is the alias of WARN.
func NameToLevel(name string) Level {
	switch strings.ToUpper(name) {
	case traceNameS:
		return LvlTrace
	case debugNameS:
		return LvlDebug
	case infoNameS:
		return LvlInfo
	case warnNameS, "WARNING":
		return LvlWarn
	case errorNameS:
		return LvlError
	case panicNameS:
		return LvlPanic
	case fatalNameS:
		return LvlFatal
	default:
		panic(fmt.Errorf("unknown level name '%s'", name))
	}
}
