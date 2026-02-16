package graph

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
)

// GraphRequest represents a request to add triples to the graph.
type GraphRequest struct {
	Triples []Triple
}

// GraphBuffer collects graph update requests and processes them in batches
// on a background goroutine. It follows the same lifecycle pattern as
// embedding.EmbeddingBuffer: Start → Enqueue → Stop.
type GraphBuffer struct {
	store Store

	queue  chan GraphRequest
	stopCh chan struct{}
	done   chan struct{}

	batchSize    int
	batchTimeout time.Duration
	logger       *zap.SugaredLogger
}

// NewGraphBuffer creates a new asynchronous graph update buffer.
func NewGraphBuffer(store Store, logger *zap.SugaredLogger) *GraphBuffer {
	return &GraphBuffer{
		store:        store,
		queue:        make(chan GraphRequest, 256),
		stopCh:       make(chan struct{}),
		done:         make(chan struct{}),
		batchSize:    64,
		batchTimeout: 2 * time.Second,
		logger:       logger,
	}
}

// Start launches the background goroutine. The WaitGroup is incremented
// so callers can wait for graceful shutdown.
func (b *GraphBuffer) Start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(b.done)
		b.run()
	}()
}

// Enqueue submits a graph update request. Non-blocking; drops if the queue is full.
func (b *GraphBuffer) Enqueue(req GraphRequest) {
	select {
	case b.queue <- req:
	default:
		b.logger.Debugw("graph queue full, dropping request", "triples", len(req.Triples))
	}
}

// Stop signals the background goroutine to drain and exit.
func (b *GraphBuffer) Stop() {
	close(b.stopCh)
	<-b.done
}

func (b *GraphBuffer) run() {
	timer := time.NewTimer(b.batchTimeout)
	defer timer.Stop()

	var batch []Triple

	flush := func() {
		if len(batch) == 0 {
			return
		}
		b.processBatch(batch)
		batch = batch[:0]
	}

	for {
		select {
		case req := <-b.queue:
			batch = append(batch, req.Triples...)
			if len(batch) >= b.batchSize {
				flush()
				timer.Reset(b.batchTimeout)
			}

		case <-timer.C:
			flush()
			timer.Reset(b.batchTimeout)

		case <-b.stopCh:
			// Drain remaining items.
			for {
				select {
				case req := <-b.queue:
					batch = append(batch, req.Triples...)
				default:
					flush()
					return
				}
			}
		}
	}
}

func (b *GraphBuffer) processBatch(batch []Triple) {
	ctx := context.Background()

	if err := b.store.AddTriples(ctx, batch); err != nil {
		b.logger.Errorw("batch graph update error", "count", len(batch), "error", err)
	}
}
