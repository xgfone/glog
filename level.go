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
