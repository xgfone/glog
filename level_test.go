package logger

import (
	"reflect"
	"testing"
)

func TestLevel(t *testing.T) {
	if LvlTrace.String() != "TRACE" {
		t.Fail()
	}
	if LvlDebug.String() != "DEBUG" {
		t.Fail()
	}
	if LvlInfo.String() != "INFO" {
		t.Fail()
	}
	if LvlWarn.String() != "WARN" {
		t.Fail()
	}
	if LvlError.String() != "ERROR" {
		t.Fail()
	}
	if LvlPanic.String() != "PANIC" {
		t.Fail()
	}
	if LvlFatal.String() != "FATAL" {
		t.Fail()
	}

	if !reflect.DeepEqual(LvlTrace.Bytes(), []byte("TRACE")) {
		t.Fail()
	}
	if !reflect.DeepEqual(LvlDebug.Bytes(), []byte("DEBUG")) {
		t.Fail()
	}
	if !reflect.DeepEqual(LvlInfo.Bytes(), []byte("INFO")) {
		t.Fail()
	}
	if !reflect.DeepEqual(LvlWarn.Bytes(), []byte("WARN")) {
		t.Fail()
	}
	if !reflect.DeepEqual(LvlError.Bytes(), []byte("ERROR")) {
		t.Fail()
	}
	if !reflect.DeepEqual(LvlPanic.Bytes(), []byte("PANIC")) {
		t.Fail()
	}
	if !reflect.DeepEqual(LvlFatal.Bytes(), []byte("FATAL")) {
		t.Fail()
	}
}

func TestNameToLevel(t *testing.T) {
	if NameToLevel("trace") != LvlTrace {
		t.Fail()
	}
	if NameToLevel("debug") != LvlDebug {
		t.Fail()
	}
	if NameToLevel("info") != LvlInfo {
		t.Fail()
	}
	if NameToLevel("warn") != LvlWarn {
		t.Fail()
	}
	if NameToLevel("warning") != LvlWarn {
		t.Fail()
	}
	if NameToLevel("error") != LvlError {
		t.Fail()
	}
	if NameToLevel("panic") != LvlPanic {
		t.Fail()
	}
	if NameToLevel("fatal") != LvlFatal {
		t.Fail()
	}
}
