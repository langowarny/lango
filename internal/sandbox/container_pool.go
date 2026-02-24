package sandbox

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ContainerPool manages a pool of pre-warmed containers for faster execution.
// It is only activated when PoolSize > 0.
type ContainerPool struct {
	runtime     ContainerRuntime
	image       string
	size        int
	idleTimeout time.Duration
	pool        chan string // buffered channel of container IDs
	mu          sync.Mutex
	closed      bool
}

// NewContainerPool creates a container pool with the specified size.
// If size is 0, the pool is effectively disabled (Acquire always returns empty).
func NewContainerPool(runtime ContainerRuntime, image string, size int, idleTimeout time.Duration) *ContainerPool {
	p := &ContainerPool{
		runtime:     runtime,
		image:       image,
		size:        size,
		idleTimeout: idleTimeout,
		pool:        make(chan string, size),
	}
	return p
}

// Acquire retrieves a pre-warmed container ID from the pool.
// If the pool is empty, it returns an empty string (caller should create on demand).
func (p *ContainerPool) Acquire(ctx context.Context) (string, error) {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return "", fmt.Errorf("pool is closed")
	}
	p.mu.Unlock()

	select {
	case id := <-p.pool:
		return id, nil
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		// Pool empty — caller should create container on demand.
		return "", nil
	}
}

// Release returns a container ID to the pool for reuse.
// If the pool is full or closed, the container is discarded.
func (p *ContainerPool) Release(containerID string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return
	}

	select {
	case p.pool <- containerID:
		// Returned to pool.
	default:
		// Pool full, discard.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = p.runtime.Cleanup(ctx, containerID)
	}
}

// Close drains the pool and cleans up all pre-warmed containers.
func (p *ContainerPool) Close() {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return
	}
	p.closed = true
	p.mu.Unlock()

	close(p.pool)
	for id := range p.pool {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_ = p.runtime.Cleanup(ctx, id)
		cancel()
	}
}

// Size returns the current number of containers in the pool.
func (p *ContainerPool) Size() int {
	return len(p.pool)
}

// Capacity returns the maximum pool capacity.
func (p *ContainerPool) Capacity() int {
	return p.size
}
