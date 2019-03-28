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

// Record stands for a log record.
type Record struct {
	// Depth is the depth of the caller.
	Depth int

	// Lvl is the level of the emitted log.
	Lvl Level

	// Msg and Args are the arguments of the emitted log.
	Msg  string
	Args []interface{}

	// Ctxs is the contexts of the Logger instance.
	Ctxs []interface{}
}
