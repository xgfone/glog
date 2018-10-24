package miss

import (
	"os"
)

func ExampleKvTextEncoder() {
	conf := EncoderConfig{IsLevel: true}
	encoder := KvTextEncoder(os.Stdout, conf)
	log := New(encoder).Cxt("name", "example", "id", 123)
	log.Info("test", "key1", "value1", "key2", "value2")

	// Output:
	// lvl=INFO name=example id=123 key1=value1 key2=value2 msg=test
}

func ExampleFmtTextEncoder() {
	conf := EncoderConfig{IsLevel: true}
	encoder := FmtTextEncoder(os.Stdout, conf)
	log := New(encoder).Cxt("kv", "text", "example")
	log.Info("test %s %s", "value1", "value2")

	// Output:
	// INFO [kv][text][example] :=>: test value1 value2
}

func ExampleLevelFilterWriter() {
	logger1 := New(KvTextEncoder(os.Stdout))
	logger1.Info("will output")

	writer := LevelFilterWriter(ERROR, os.Stdout)
	logger2 := New(KvTextEncoder(writer))
	logger2.Info("won't output")

	// Output:
	// will output
}
