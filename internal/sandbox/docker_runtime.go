package sandbox

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// DockerRuntime executes tools inside Docker containers.
type DockerRuntime struct {
	cli *client.Client
}

// NewDockerRuntime creates a DockerRuntime using the default Docker client.
func NewDockerRuntime() (*DockerRuntime, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("create docker client: %w", err)
	}
	return &DockerRuntime{cli: cli}, nil
}

// Run executes a tool inside a Docker container, communicating via stdin/stdout JSON.
func (r *DockerRuntime) Run(ctx context.Context, cfg ContainerConfig) (*ExecutionResult, error) {
	if cfg.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cfg.Timeout)
		defer cancel()
	}

	// Prepare the execution request.
	req := ExecutionRequest{
		Version:  1,
		ToolName: cfg.ToolName,
		Params:   cfg.Params,
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal execution request: %w", err)
	}

	// Container configuration.
	containerCfg := &container.Config{
		Image:        cfg.Image,
		Cmd:          []string{"--sandbox-worker"},
		OpenStdin:    true,
		StdinOnce:    true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Labels: map[string]string{
			"lango.sandbox": "true",
			"lango.tool":    cfg.ToolName,
		},
	}

	hostCfg := &container.HostConfig{
		NetworkMode: container.NetworkMode(cfg.NetworkMode),
		Resources: container.Resources{
			Memory:   cfg.MemoryLimitMB * 1024 * 1024,
			CPUQuota: cfg.CPUQuotaUS,
		},
		ReadonlyRootfs: cfg.ReadOnlyRootfs,
		Tmpfs: map[string]string{
			"/tmp": "",
		},
	}

	// Create the container.
	resp, err := r.cli.ContainerCreate(ctx, containerCfg, hostCfg, nil, nil, "")
	if err != nil {
		return nil, fmt.Errorf("create container: %w", err)
	}
	containerID := resp.ID

	// Always remove the container when done.
	defer func() {
		removeCtx, removeCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer removeCancel()
		_ = r.cli.ContainerRemove(removeCtx, containerID, container.RemoveOptions{Force: true})
	}()

	// Attach to container for stdin/stdout hijacking.
	attachResp, err := r.cli.ContainerAttach(ctx, containerID, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return nil, fmt.Errorf("attach container: %w", err)
	}
	defer attachResp.Close()

	// Start the container.
	if err := r.cli.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
		return nil, fmt.Errorf("start container: %w", err)
	}

	// Write JSON request to stdin, then close.
	if _, err := attachResp.Conn.Write(reqBytes); err != nil {
		return nil, fmt.Errorf("write to container stdin: %w", err)
	}
	if err := attachResp.CloseWrite(); err != nil {
		return nil, fmt.Errorf("close container stdin: %w", err)
	}

	// Read stdout.
	var stdout bytes.Buffer
	if _, err := io.Copy(&stdout, attachResp.Reader); err != nil {
		// Ignore read errors if context was cancelled (timeout).
		if ctx.Err() != nil {
			return nil, ErrContainerTimeout
		}
		return nil, fmt.Errorf("read container stdout: %w", err)
	}

	// Wait for the container to finish.
	waitCh, errCh := r.cli.ContainerWait(ctx, containerID, container.WaitConditionNotRunning)
	select {
	case waitResp := <-waitCh:
		// Exit code 137 indicates OOM kill (128 + SIGKILL=9).
		if waitResp.StatusCode == 137 {
			return nil, ErrContainerOOM
		}
	case err := <-errCh:
		if ctx.Err() != nil {
			return nil, ErrContainerTimeout
		}
		return nil, fmt.Errorf("wait for container: %w", err)
	case <-ctx.Done():
		return nil, ErrContainerTimeout
	}

	// Parse result from stdout.
	// Docker multiplexes stdout/stderr with an 8-byte header per frame.
	// Try to parse the raw output first; if it fails, strip headers.
	rawOutput := stdout.Bytes()
	var result ExecutionResult
	if err := json.Unmarshal(rawOutput, &result); err != nil {
		// Try stripping Docker stream headers (8-byte prefix per frame).
		stripped := stripDockerStreamHeaders(rawOutput)
		if jsonErr := json.Unmarshal(stripped, &result); jsonErr != nil {
			return nil, fmt.Errorf("unmarshal container result: %w (raw: %s)", err, string(rawOutput))
		}
	}

	if result.Error != "" {
		return &result, fmt.Errorf("tool %q: %s", cfg.ToolName, result.Error)
	}

	return &result, nil
}

// stripDockerStreamHeaders removes Docker multiplexed stream headers.
// Each frame has: [type(1)][padding(3)][size(4)][payload(size)].
func stripDockerStreamHeaders(data []byte) []byte {
	var out bytes.Buffer
	for len(data) >= 8 {
		size := int(data[4])<<24 | int(data[5])<<16 | int(data[6])<<8 | int(data[7])
		data = data[8:]
		if size > len(data) {
			size = len(data)
		}
		out.Write(data[:size])
		data = data[size:]
	}
	return out.Bytes()
}

// Cleanup removes orphaned sandbox containers with the "lango.sandbox=true" label.
func (r *DockerRuntime) Cleanup(ctx context.Context, _ string) error {
	containers, err := r.cli.ContainerList(ctx, container.ListOptions{
		All: true,
	})
	if err != nil {
		return fmt.Errorf("list containers: %w", err)
	}

	var removed int
	for _, c := range containers {
		if c.Labels["lango.sandbox"] == "true" {
			if err := r.cli.ContainerRemove(ctx, c.ID, container.RemoveOptions{Force: true}); err != nil {
				continue
			}
			removed++
		}
	}
	return nil
}

// IsAvailable checks if Docker daemon is reachable.
func (r *DockerRuntime) IsAvailable(ctx context.Context) bool {
	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	_, err := r.cli.Ping(pingCtx)
	return err == nil
}

// Name returns the runtime name.
func (r *DockerRuntime) Name() string {
	return "docker"
}
