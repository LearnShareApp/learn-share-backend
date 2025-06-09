package workerpool

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"sync"
)

type WorkerPool[T any] struct {
	maxParallelOperation int
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
	var mu sync.Mutex

	eg, ctx := errgroup.WithContext(ctx)

	// limited count of parallel requests
	sem := make(chan struct{}, w.maxParallelOperation)

	for id := range unfilledMap {
		idCopy := id // local copy for gorutine (shadow variable)
		eg.Go(func() error {
			sem <- struct{}{}        // book slot
			defer func() { <-sem }() // unbook slot

			dataForFilling, err := fillerFunc(ctx, idCopy)
			if err != nil {
				return err
			}

			mu.Lock()
			defer mu.Unlock()
			unfilledMap[idCopy] = dataForFilling

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("worker pool, FillMap: failed to get data by fillerFunc: %w", err)
	}

	return nil
}
