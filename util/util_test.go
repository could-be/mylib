package util

import (
	"errors"
	"log"
	"os"
	"testing"
)

func TestPanicOnError(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	_ = os.Setenv(ModeFlag, "test")
	// 为空
	PanicOnError(nil, "test %s", "test")
	PanicOnError(nil, "test2")

	// 不为空
	// PanicOnError(errors.New("TEST NOT NIL"), "test for (%s, %d)", "test", 0)
	PanicOnError(errors.New("TEST NOT NIL"), "test empty params")

}
