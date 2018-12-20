package miss

import (
	"bytes"
	"testing"
)

func TestLoggerWithoutError(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	errlog := New(FmtTextEncoder(buf)).Cxt(Caller())
	logger := ToLoggerWithoutError(errlog)

	errlog.Info("hello, %s", "abc")
	logger.Info("hello, %s", "xyz")
	if buf.String() != "[noerrorlog_test.go:13] :=>: hello, abc\n[noerrorlog_test.go:14] :=>: hello, xyz\n" {
		t.Error(buf.String())
		t.Fail()
	}

	if ToLogger(logger).GetDepth() != errlog.GetDepth() {
		t.Fail()
	}
}
