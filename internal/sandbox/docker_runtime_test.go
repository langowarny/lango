package sandbox

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDockerRuntime_Name(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Docker integration test in short mode")
	}
	dr, err := NewDockerRuntime()
	if err != nil {
		t.Skipf("Docker client unavailable: %v", err)
	}
	assert.Equal(t, "docker", dr.Name())
}

func TestStripDockerStreamHeaders(t *testing.T) {
	tests := []struct {
		give string
		want string
	}{
		{
			give: "already plain JSON",
			want: "already plain JSON",
		},
	}

	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			// Build a Docker-style frame: [type=1][0,0,0][size_be32][payload].
			payload := []byte(tt.want)
			frame := make([]byte, 8+len(payload))
			frame[0] = 1 // stdout
			frame[4] = byte(len(payload) >> 24)
			frame[5] = byte(len(payload) >> 16)
			frame[6] = byte(len(payload) >> 8)
			frame[7] = byte(len(payload))
			copy(frame[8:], payload)

			result := stripDockerStreamHeaders(frame)
			assert.True(t, bytes.Equal(payload, result))
		})
	}
}

func TestStripDockerStreamHeaders_MultipleFrames(t *testing.T) {
	part1 := []byte(`{"output":`)
	part2 := []byte(`{"ok":true}}`)

	var buf bytes.Buffer
	// Frame 1
	frame1 := make([]byte, 8+len(part1))
	frame1[0] = 1
	frame1[7] = byte(len(part1))
	copy(frame1[8:], part1)
	buf.Write(frame1)

	// Frame 2
	frame2 := make([]byte, 8+len(part2))
	frame2[0] = 1
	frame2[7] = byte(len(part2))
	copy(frame2[8:], part2)
	buf.Write(frame2)

	result := stripDockerStreamHeaders(buf.Bytes())
	expected := append(part1, part2...)
	assert.Equal(t, expected, result)
}

func TestStripDockerStreamHeaders_EmptyInput(t *testing.T) {
	result := stripDockerStreamHeaders(nil)
	assert.Empty(t, result)
}
