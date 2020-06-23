package lru

import (
	"runtime"
	"testing"
)

// LRU 最近最少使用，内存淘汰算法
func TestExample(t *testing.T) {
	l := NewLRU(10)
	l.Assign(1, 1)
	l.Print()
	// 强制性gc, 设置回收gc时的回调函数
	runtime.GC()
	runtime.SetFinalizer(l, func(v interface{}) {

	})
}
