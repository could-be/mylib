package concurrency

import (
	"context"
	"errors"
	"fmt"
)

// 有输入源的是有效额
type ConcurrencyEntity struct {
	workers     map[string]*worker
	seed        string
	err         chan error
	validWorker map[string]bool

	terminator string                          // pipeline 的最终一个
	finalizer  func(in chan interface{}) error // finalizer 阻塞读取管道最后数据，直到engine所有任务结束
	// pipeline []string
}

func New() *ConcurrencyEntity {

	return &ConcurrencyEntity{
		workers:     map[string]*worker{},
		err:         make(chan error),
		validWorker: make(map[string]bool),
		// pipeline: make([]string, 0, 10),
	}
}

// 最后停止，依赖管道最后一个
// BUG: 有一定几率，直接返回，不运行
// 添加任务太快了，wait的时候，添加任务都还没有添加寝取呢
// wait 执行在go add 之前了
func (s *ConcurrencyEntity) Run() error {
	// Check Run之前，管道最终一环必须设置
	if s.terminator == "" {
		return errors.New("invalid terminator")
	}

	select {
	default:
	case e := <-s.err:
		return e
	}

	if err := s.checkPipeline(); err != nil {
		return err
	}
	s.waitClose()

	return s.finalizer(s.result())
}

func (s *ConcurrencyEntity) SetFinalizer(workerName string, f func(in chan interface{}) error) {
	s.terminator = workerName
	s.finalizer = f
}

// 在Run方法之后调用，Run方法检查Terminator合法性
// TODO: 需要检查的
func (s *ConcurrencyEntity) result() chan interface{} {
	return s.workers[s.terminator].C()
}

// TODO: 打印pipeline地图
func (s *ConcurrencyEntity) PrintPipeline() string {
	return ""
}

func (s *ConcurrencyEntity) checkPipeline() error {
	for pipe, valid := range s.validWorker {
		if pipe != s.seed && !valid {
			return fmt.Errorf("pipeline worker(%s) has no source pipe", pipe)
		}
	}
	return nil
}

// func (s *ConcurrencyEntity) AddPageSeed(ctx context.Context,
// )

// 暂时只支持一个种子
func (s *ConcurrencyEntity) AddPageSeed(ctx context.Context,
	pageSize, total int, seed string, f func(offset, page int, out chan interface{}) error) {
	if _, ok := s.workers[seed]; ok {
		s.err <- fmt.Errorf("seed alreead exists")
	}
	s.seed = seed
	w := s.Worker(ctx, seed)
	page := (total + pageSize) / pageSize

	for i := 0; i < pageSize; i++ {
		i := i
		w.Add(func() error {
			return f(i*page, page, w.C())
		})
	}
}

func (s *ConcurrencyEntity) AddSingleSeed(ctx context.Context,
	seed string, f func(out chan interface{}) error) {
	if _, ok := s.workers[seed]; ok {
		s.err <- fmt.Errorf("seed alreead exists")
	}
	s.seed = seed
	w := s.Worker(ctx, seed)

	w.Add(func() error {
		return f(w.C())
	})
}

func (s *ConcurrencyEntity) Worker(ctx context.Context, pipeName string) *worker {
	g, ok := s.workers[pipeName]
	if ok {
		return g
	}

	g = newWorker(ctx)
	s.workers[pipeName] = g
	s.validWorker[pipeName] = false
	return g
}

func (s *ConcurrencyEntity) withValid(pipeName string) {
	s.validWorker[pipeName] = true
}

func (s *ConcurrencyEntity) EmptyWait(in chan interface{}) error {
	for _ = range in {
	}
	return nil
}

func (s *ConcurrencyEntity) Add(ctx context.Context, n int, from, to string, f func(in, out chan interface{}) error) {
	fromW := s.Worker(ctx, from)
	toW := s.Worker(ctx, to)
	s.withValid(to)

	// 第一个函数式为了，lasy使用内部结构的字段; 调用方决定使用那个入参
	// 第二个是闭包，是为了使用当前上下文环境; 昂前上下文包含了 当前作用范围和上一层的函数参数
	for i := 0; i < n; i++ {
		toW.Add(func() error {
			return f(fromW.C(), toW.C())
		})
	}
}

func (s *ConcurrencyEntity) waitClose() {
	for _, v := range s.workers {
		v.GoWaitClose()
	}
}
