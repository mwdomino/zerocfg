package zerocfg

import (
	"runtime"
	"strings"
)

const (
	configPackage = "github.com/chaindead/zerocfg."
	initSuffix    = ".init"
)

func findCaller() string {
	var nextToZfg bool
	for i := 2; ; i++ {
		pc, _, _, _ := runtime.Caller(i)
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			return ""
		}

		name := fn.Name()
		if !nextToZfg && !strings.HasPrefix(name, configPackage) {
			nextToZfg = true
		}

		if nextToZfg && strings.HasSuffix(name, initSuffix) {
			return strings.TrimSuffix(name, initSuffix)
		}
	}
}
