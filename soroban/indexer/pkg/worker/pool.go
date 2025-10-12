package worker

import (
	"context"
	"sync"
)

// Task represents a unit of work to be processed
type Task func(ctx context.Context) error

// Result represents the outcome of a task execution
type Result struct {
	Error error
	Index int // Original index in the task list
}

// Pool manages a pool of workers that execute tasks concurrently
type Pool struct {
	maxWorkers int
	taskQueue  chan Task
	results    []Result
	resultsMu  sync.Mutex
	wg         sync.WaitGroup
}

// NewPool creates a new worker pool with the specified number of workers
func NewPool(maxWorkers int) *Pool {
	if maxWorkers <= 0 {
		maxWorkers = 1
	}

	return &Pool{
		maxWorkers: maxWorkers,
		taskQueue:  make(chan Task, maxWorkers*2), // Buffer to prevent blocking
		results:    make([]Result, 0),
	}
}

// Start begins the worker pool
func (p *Pool) Start(ctx context.Context) {
	for i := 0; i < p.maxWorkers; i++ {
		p.wg.Add(1)
		go p.worker(ctx)
	}
}

// worker processes tasks from the task queue
func (p *Pool) worker(ctx context.Context) {
	defer p.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-p.taskQueue:
			if !ok {
				return
			}
			// Execute the task
			err := task(ctx)

			// Store result safely
			p.resultsMu.Lock()
			p.results = append(p.results, Result{Error: err})
			p.resultsMu.Unlock()
		}
	}
}

// Submit adds a task to the worker pool
func (p *Pool) Submit(task Task) {
	p.taskQueue <- task
}

// Wait closes the task queue and waits for all workers to finish
// Returns all results
func (p *Pool) Wait() []Result {
	close(p.taskQueue)
	p.wg.Wait()

	p.resultsMu.Lock()
	defer p.resultsMu.Unlock()

	return p.results
}

// Execute runs multiple tasks concurrently and collects results
// This is a convenience method that creates a pool, executes tasks, and waits
func Execute(ctx context.Context, maxWorkers int, tasks []Task) []Result {
	pool := NewPool(maxWorkers)
	pool.Start(ctx)

	for _, task := range tasks {
		pool.Submit(task)
	}

	return pool.Wait()
}
