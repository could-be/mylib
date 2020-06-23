package util

import (
	"os"

	"github.com/couldbe/fruit/infrastructure/ilog"
)

// TODO: 替换自己的日志框架.
func PanicOnError(err error, format string, v ...interface{}) {
	switch err {
	default:
		v = append(v, err)
		ilog.ErrorfDepth(1, format, v...)
		panic(err)
	case nil:
		switch {
		case v != nil && len(v) == 1 && v[0] == nil:
			v[0] = "success"
		default:
			format += "%s"
			v = append(v, "success")
		}
		ilog.V(4).InfofDepth(1, format, v...)
	}

}

var runLevel RunLevel

var RunMode = "RUN_MODE"

func init() {
	switch os.Getenv(RunMode) {
	case TestMode:
		runLevel = 1
	case DebugMode:
		runLevel = 2
	case InitMode:
		runLevel = 4
	}
}

type RunLevel int

const (
	ReleaseLevel RunLevel = iota
	TestLevel
	DebugLevel
	InitLevel
)

const (
	ReleaseMode = "release"
	TestMode    = "test"
	DebugMode   = "debug"
	InitMode    = "init"
)

func (l RunLevel) String() string {
	switch l {
	case 1:
		return TestMode
	case 2:
		return DebugMode
	case 6:
		return InitMode
	default:
		return ReleaseMode
	}
}

func Debug() bool {
	return runLevel >= DebugLevel
}

func Test() bool {
	return runLevel >= TestLevel
}

func Init() bool {
	return runLevel >= InitLevel
}
func Release() bool {
	return runLevel >= ReleaseLevel
}
