package concurrency

import "context"

type Concurrency interface {
	Run() error
	SetFinalizer(workerName string, f func(in chan interface{}) error)
	AddPageSeed(ctx context.Context,
		n, total int, seed string, f func(offset, page int, out chan interface{}) error)
	AddSingleSeed(ctx context.Context,
		seed string, f func(out chan interface{}) error)
	EmptyWait(in chan interface{}) error
	Add(ctx context.Context, n int, from, to string, f func(in, out chan interface{}) error)
}
