// Copyright 2019 xgfone
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

	"github.com/xgfone/logger/utils"
)

// Re-export some errors.
var (
	ErrType = utils.ErrType
)

// MultiError represents more than one error.
type MultiError struct {
	errs []error
}

func (m MultiError) Error() string {
	switch len(m.errs) {
	case 0:
		return ""
	case 1:
		return m.errs[0].Error()
	default:
		s := m.errs[0].Error()
		for _, e := range m.errs[1:] {
			if e != nil {
				s = fmt.Sprintf("%s > %s", s, e.Error())
			}
		}
		return s
	}
}

// Errors returns a list of errors.
func (m MultiError) Errors() []error {
	return m.errs
}
