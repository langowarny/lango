package sandbox

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContainerPool_AcquireRelease(t *testing.T) {
	rt := &mockRuntime{name: "mock", available: true}
	pool := NewContainerPool(rt, "test:latest", 3, 5*time.Minute)
	defer pool.Close()

	// Pool starts empty.
	assert.Equal(t, 0, pool.Size())
	assert.Equal(t, 3, pool.Capacity())

	// Acquire from empty pool returns empty string (no block).
	id, err := pool.Acquire(context.Background())
	require.NoError(t, err)
	assert.Empty(t, id)

	// Release a container ID.
	pool.Release("container-1")
	assert.Equal(t, 1, pool.Size())

	// Acquire retrieves it.
	id, err = pool.Acquire(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "container-1", id)
	assert.Equal(t, 0, pool.Size())
}

func TestContainerPool_ReleaseFullPool(t *testing.T) {
	rt := &mockRuntime{name: "mock", available: true}
	pool := NewContainerPool(rt, "test:latest", 2, 5*time.Minute)
	defer pool.Close()

	pool.Release("c-1")
	pool.Release("c-2")
	assert.Equal(t, 2, pool.Size())

	// Releasing when full should discard.
	pool.Release("c-3")
	assert.Equal(t, 2, pool.Size())
}

func TestContainerPool_Close(t *testing.T) {
	rt := &mockRuntime{name: "mock", available: true}
	pool := NewContainerPool(rt, "test:latest", 3, 5*time.Minute)

	pool.Release("c-1")
	pool.Release("c-2")
	pool.Close()

	// Acquire after close returns error.
	_, err := pool.Acquire(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "pool is closed")

	// Release after close is no-op.
	pool.Release("c-3") // should not panic
}

func TestContainerPool_DoubleClose(t *testing.T) {
	rt := &mockRuntime{name: "mock", available: true}
	pool := NewContainerPool(rt, "test:latest", 2, 5*time.Minute)

	pool.Close()
	pool.Close() // should not panic
}

func TestContainerPool_AcquireContextCancelled(t *testing.T) {
	rt := &mockRuntime{name: "mock", available: true}
	pool := NewContainerPool(rt, "test:latest", 3, 5*time.Minute)
	defer pool.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// With cancelled context and empty pool, Acquire may return either
	// empty string (default case) or context error (ctx.Done case).
	id, err := pool.Acquire(ctx)
	if err != nil {
		assert.ErrorIs(t, err, context.Canceled)
	}
	assert.Empty(t, id)
}

func TestContainerPool_FIFO(t *testing.T) {
	rt := &mockRuntime{name: "mock", available: true}
	pool := NewContainerPool(rt, "test:latest", 5, 5*time.Minute)
	defer pool.Close()

	pool.Release("c-1")
	pool.Release("c-2")
	pool.Release("c-3")

	id1, _ := pool.Acquire(context.Background())
	id2, _ := pool.Acquire(context.Background())
	id3, _ := pool.Acquire(context.Background())

	assert.Equal(t, "c-1", id1)
	assert.Equal(t, "c-2", id2)
	assert.Equal(t, "c-3", id3)
}
