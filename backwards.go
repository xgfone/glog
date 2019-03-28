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
	"github.com/xgfone/logger/utils"
)

// Re-export some errors.
var (
	ErrType = utils.ErrType
)

// Re-export some functions for backward compatibility.
var (
	Range = utils.Range

	ToBytes     = utils.ToBytes
	ToBytesErr  = utils.ToBytesErr
	ToString    = utils.ToString
	ToStringErr = utils.ToStringErr

	WriteString     = utils.WriteString
	WriteIntoBuffer = utils.WriteIntoBuffer

	MarshalJSON   = utils.MarshalJSON
	MarshalKvJSON = utils.MarshalKvJSON
)

type (
	// Byter is re-exported for backward compatibility.
	Byter = utils.Byter

	// MarshalText is re-exported for backward compatibility.
	MarshalText = utils.MarshalText

	// StringWriter is re-exported for backward compatibility.
	StringWriter = utils.StringWriter
)
