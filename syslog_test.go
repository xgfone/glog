// +build !windows,!plan9

package miss

import (
	"fmt"
	"log/syslog"
)

func ExampleSyslogWriter() {
	conf := EncoderConfig{IsLevel: true}
	out, c := Must.SyslogWriter(syslog.LOG_DEBUG, "testsyslog")
	defer c.Close()
	log := New(FmtTextEncoder(out, conf))
	if err := log.Info("test %s %s", "syslog", "writer"); err != nil {
		fmt.Printf("Error: %s\n", err)
	}

	// Output:
	//
}
