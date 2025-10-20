package workerpool

import (
	"context"
	"fmt"
	"sync"
)

type WorkerPool[T any] struct {
	maxParallelOperation int
}

type job[T any] struct {
	id int
}

type result[T any] struct {
	id   int
	data *T
	err  error
}

func NewWorkerPool[T any](maxParallelOperation int) *WorkerPool[T] {
	if maxParallelOperation <= 0 {
		maxParallelOperation = 1
	}
	return &WorkerPool[T]{
		maxParallelOperation: maxParallelOperation,
	}
}

func (w *WorkerPool[T]) FillMap(
	ctx context.Context,
	unfilledMap map[int]*T,
	fillerFunc func(ctx context.Context, id int) (*T, error)) error {

	if len(unfilledMap) == 0 {
		return nil
	}

	jobs := make(chan job[T])
	results := make(chan result[T], w.maxParallelOperation)

	// Запускаем воркеров
	var wg sync.WaitGroup
	for i := 0; i < w.maxParallelOperation; i++ {
		wg.Add(1)
		go w.worker(ctx, jobs, results, fillerFunc, &wg)
	}

	// Горутина для отправки задач
	go func() {
		defer close(jobs)
		for id := range unfilledMap {
			select {
			case jobs <- job[T]{id: id}:
			case <-ctx.Done():
				return
			}
		}
	}()

	// Горутина для закрытия канала результатов после завершения всех воркеров
	go func() {
		wg.Wait()
		close(results)
	}()

	return w.collectResults(ctx, unfilledMap, results)
}

func (w *WorkerPool[T]) worker(
	ctx context.Context,
	jobs <-chan job[T],
	results chan<- result[T],
	fillerFunc func(ctx context.Context, id int) (*T, error),
	wg *sync.WaitGroup) {

	defer wg.Done()

	for {
		select {
		case j, ok := <-jobs:
			if !ok {
				return
			}

			select {
			case <-ctx.Done():
				return
			default:
			}

			data, err := fillerFunc(ctx, j.id)

			select {
			case results <- result[T]{id: j.id, data: data, err: err}:
			case <-ctx.Done():
				return
			}

		case <-ctx.Done():
			return
		}
	}
}

func (w *WorkerPool[T]) collectResults(
	ctx context.Context,
	unfilledMap map[int]*T,
	results <-chan result[T]) error {

	processed := 0
	expectedCount := len(unfilledMap)

	for {
		select {
		case res, ok := <-results:
			if !ok {
				if processed < expectedCount {
					return fmt.Errorf("worker pool: not all results received, got %d out of %d", processed, expectedCount)
				}
				return nil
			}

			processed++

			if res.err != nil {
				return fmt.Errorf("worker pool, FillMap: failed to get data for id %d: %w", res.id, res.err)
			}

			unfilledMap[res.id] = res.data

			if processed == expectedCount {
				return nil
			}

		case <-ctx.Done():
			return fmt.Errorf("worker pool: context cancelled after processing %d out of %d items: %w", processed, expectedCount, ctx.Err())
		}
	}
}
