package embedding

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
)

// EmbedRequest represents a request to embed and store a text.
type EmbedRequest struct {
	ID         string
	Collection string
	Content    string
	Metadata   map[string]string
}

// EmbeddingBuffer collects embed requests and processes them in batches
// on a background goroutine. It follows the same lifecycle pattern as
// memory.Buffer: Start → Enqueue → Stop.
type EmbeddingBuffer struct {
	provider EmbeddingProvider
	store    VectorStore

	queue  chan EmbedRequest
	stopCh chan struct{}
	done   chan struct{}

	batchSize    int
	batchTimeout time.Duration
	logger       *zap.SugaredLogger
}

// NewEmbeddingBuffer creates a new asynchronous embedding buffer.
func NewEmbeddingBuffer(
	provider EmbeddingProvider,
	store VectorStore,
	logger *zap.SugaredLogger,
) *EmbeddingBuffer {
	return &EmbeddingBuffer{
		provider:     provider,
		store:        store,
		queue:        make(chan EmbedRequest, 256),
		stopCh:       make(chan struct{}),
		done:         make(chan struct{}),
		batchSize:    32,
		batchTimeout: 2 * time.Second,
		logger:       logger,
	}
}

// Start launches the background goroutine. The WaitGroup is incremented
// so callers can wait for graceful shutdown.
func (b *EmbeddingBuffer) Start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(b.done)
		b.run()
	}()
}

// Enqueue submits an embed request. Non-blocking; drops if the queue is full.
func (b *EmbeddingBuffer) Enqueue(req EmbedRequest) {
	select {
	case b.queue <- req:
	default:
		b.logger.Debugw("embedding queue full, dropping request", "id", req.ID, "collection", req.Collection)
	}
}

// Stop signals the background goroutine to drain and exit.
func (b *EmbeddingBuffer) Stop() {
	close(b.stopCh)
	<-b.done
}

func (b *EmbeddingBuffer) run() {
	timer := time.NewTimer(b.batchTimeout)
	defer timer.Stop()

	var batch []EmbedRequest

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
			batch = append(batch, req)
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
					batch = append(batch, req)
				default:
					flush()
					return
				}
			}
		}
	}
}

func (b *EmbeddingBuffer) processBatch(batch []EmbedRequest) {
	ctx := context.Background()

	texts := make([]string, len(batch))
	for i, r := range batch {
		texts[i] = r.Content
	}

	embeddings, err := b.provider.Embed(ctx, texts)
	if err != nil {
		b.logger.Errorw("batch embedding failed", "count", len(batch), "error", err)
		return
	}

	if len(embeddings) != len(batch) {
		b.logger.Errorw("embedding count mismatch", "expected", len(batch), "got", len(embeddings))
		return
	}

	records := make([]VectorRecord, len(batch))
	for i, r := range batch {
		records[i] = VectorRecord{
			ID:         r.ID,
			Collection: r.Collection,
			Embedding:  embeddings[i],
			Metadata:   r.Metadata,
		}
	}

	if err := b.store.Upsert(ctx, records); err != nil {
		b.logger.Errorw("batch upsert failed", "count", len(records), "error", err)
	}
}
