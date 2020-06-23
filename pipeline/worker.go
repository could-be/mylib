package concurrency

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"
)

// 为 go {for range} 设计的结构
// 解决多个并发写channel, pipeline 顺序、安全关闭问题

// TODO: 还可以尝试下另一种方法
//  for range {go}
type worker struct {
	g    *errgroup.Group
	ch   chan interface{}
	once *sync.Once
}

func newWorker(ctx context.Context) *worker {
	g, _ := errgroup.WithContext(ctx)

	return &worker{
		g:    g,
		ch:   make(chan interface{}),
		once: &sync.Once{},
	}
}

func (w *worker) Add(f func() error) {
	w.g.Go(f)
}

func (w *worker) C() chan interface{} {
	return w.ch
}

// 外层waite使用位置不当，可能会导致问题
func (w *worker) GoWaitClose() {
	go w.once.Do(func() {
		w.g.Wait()
		close(w.ch)
	})
}
