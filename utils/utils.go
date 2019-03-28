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

// Package utils supplies some utility functions.
package utils

import "fmt"

// Range returns a integer range between start and stop, which progressively
// increase or descrease by step.
//
// If step is positive, r[i] = start + step*i when i>0 and r[i]<stop.
//
// If step is negative, r[i] = start + step*i but when i>0 and r[i]>stop.
//
// If step is 0, it will panic.
func Range(start, stop, step int) (r []int) {
	if step > 0 {
		for start < stop {
			r = append(r, start)
			start += step
		}
		return
	} else if step < 0 {
		for start > stop {
			r = append(r, start)
			start += step
		}
		return
	}

	panic(fmt.Errorf("The step must not be 0"))
}
