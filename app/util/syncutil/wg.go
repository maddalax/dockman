package syncutil

import "sync"

type WaitGroupWithConcurrency struct {
	concurrency int
	wg          sync.WaitGroup
	semaphore   chan struct{}
}

func NewWaitGroupWithConcurrency(concurrency int) *WaitGroupWithConcurrency {
	return &WaitGroupWithConcurrency{
		concurrency: concurrency,
		wg:          sync.WaitGroup{},
		semaphore:   make(chan struct{}, concurrency)}
}

func (wg *WaitGroupWithConcurrency) Add() {
	wg.semaphore <- struct{}{}
	wg.wg.Add(1)
}

func (wg *WaitGroupWithConcurrency) Done() {
	wg.wg.Done()
	<-wg.semaphore
}

func (wg *WaitGroupWithConcurrency) Wait() {
	wg.wg.Wait()
}
