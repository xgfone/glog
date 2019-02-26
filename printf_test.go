package logger

import (
	"bytes"
	"testing"
)

func TestPrintfer(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	log := GetGlobalLogger().Encoder(FmtTextEncoder(buf))
	plog := ToPrintfer(log)
	plog.Printf("printf: %s", "abc")

	if buf.String() != "printf: abc\n" {
		t.Error(buf.String())
	}
}
