package workerpool

import (
	"sync"
)

type WorkerPool struct {
	maxWorker int
	done      chan struct{}
	errc      chan error
	onceC     *sync.Once
	wg        *sync.WaitGroup
}

func NewPool(maxWorker int) *WorkerPool {
	return &WorkerPool{
		maxWorker: maxWorker,
		done:      make(chan struct{}),
		errc:      make(chan error, 1),
		onceC:     &sync.Once{},
		wg:        &sync.WaitGroup{},
	}
}

func (wp *WorkerPool) Run(queue chan func() error) <-chan error {
	for i := 0; i < wp.maxWorker; i++ {
		wp.wg.Add(1)
		go wp.worker(queue)
	}

	go func() {
		wp.wg.Wait()
		close(wp.errc)
	}()

	return wp.errc
}

func (wp *WorkerPool) Stop() {
	// using Once to prevent panic if multiple stop called.
	//  (closing a clseod channel will panic)
	wp.onceC.Do(func() {
		close(wp.done)
	})
}

func (wp *WorkerPool) worker(queue chan func() error) {
	defer wp.wg.Done()

	for {
		select {
		case task, ok := <-queue:
			if !ok {
				return
			}

			wp.errc <- task()
		case <-wp.done:
			return
		}
	}
}
