package p2p

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/langoai/lango/internal/bootstrap"
	"github.com/langoai/lango/internal/sandbox"
)

func newSandboxCmd(bootLoader func() (*bootstrap.Result, error)) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sandbox",
		Short: "Manage P2P tool execution sandbox",
		Long:  "Inspect sandbox status, run smoke tests, and clean up orphaned containers.",
	}

	cmd.AddCommand(newSandboxStatusCmd(bootLoader))
	cmd.AddCommand(newSandboxTestCmd(bootLoader))
	cmd.AddCommand(newSandboxCleanupCmd(bootLoader))

	return cmd
}

func newSandboxStatusCmd(bootLoader func() (*bootstrap.Result, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show sandbox runtime status",
		RunE: func(cmd *cobra.Command, args []string) error {
			boot, err := bootLoader()
			if err != nil {
				return err
			}

			cfg := boot.Config
			if !cfg.P2P.ToolIsolation.Enabled {
				fmt.Println("Tool isolation: disabled")
				return nil
			}

			fmt.Println("Tool isolation: enabled")
			fmt.Printf("  Timeout per tool: %v\n", cfg.P2P.ToolIsolation.TimeoutPerTool)
			fmt.Printf("  Max memory (MB):  %d\n", cfg.P2P.ToolIsolation.MaxMemoryMB)

			if !cfg.P2P.ToolIsolation.Container.Enabled {
				fmt.Println("  Container mode:   disabled (subprocess fallback)")
				return nil
			}

			fmt.Println("  Container mode:   enabled")
			fmt.Printf("  Runtime config:   %s\n", cfg.P2P.ToolIsolation.Container.Runtime)
			fmt.Printf("  Image:            %s\n", cfg.P2P.ToolIsolation.Container.Image)
			fmt.Printf("  Network mode:     %s\n", cfg.P2P.ToolIsolation.Container.NetworkMode)

			// Probe actual runtime availability.
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			sbxCfg := sandbox.Config{
				Enabled:        true,
				TimeoutPerTool: cfg.P2P.ToolIsolation.TimeoutPerTool,
				MaxMemoryMB:    cfg.P2P.ToolIsolation.MaxMemoryMB,
			}
			exec, err := sandbox.NewContainerExecutor(sbxCfg, cfg.P2P.ToolIsolation.Container)
			if err != nil {
				fmt.Printf("  Active runtime:   unavailable (%v)\n", err)
				return nil
			}
			_ = ctx
			fmt.Printf("  Active runtime:   %s\n", exec.RuntimeName())
			fmt.Printf("  Pool size:        %d\n", cfg.P2P.ToolIsolation.Container.PoolSize)

			return nil
		},
	}
}

func newSandboxTestCmd(bootLoader func() (*bootstrap.Result, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "test",
		Short: "Run a sandbox smoke test",
		Long:  "Execute a simple echo tool through the sandbox to verify it works.",
		RunE: func(cmd *cobra.Command, args []string) error {
			boot, err := bootLoader()
			if err != nil {
				return err
			}

			cfg := boot.Config
			if !cfg.P2P.ToolIsolation.Enabled {
				return fmt.Errorf("tool isolation is not enabled (set p2p.toolIsolation.enabled = true)")
			}

			sbxCfg := sandbox.Config{
				Enabled:        true,
				TimeoutPerTool: cfg.P2P.ToolIsolation.TimeoutPerTool,
				MaxMemoryMB:    cfg.P2P.ToolIsolation.MaxMemoryMB,
			}

			var exec sandbox.Executor
			if cfg.P2P.ToolIsolation.Container.Enabled {
				containerExec, cErr := sandbox.NewContainerExecutor(sbxCfg, cfg.P2P.ToolIsolation.Container)
				if cErr != nil {
					fmt.Printf("Container sandbox unavailable, using subprocess: %v\n", cErr)
					exec = sandbox.NewSubprocessExecutor(sbxCfg)
				} else {
					fmt.Printf("Using container runtime: %s\n", containerExec.RuntimeName())
					exec = containerExec
				}
			} else {
				fmt.Println("Using subprocess sandbox")
				exec = sandbox.NewSubprocessExecutor(sbxCfg)
			}

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			params := map[string]interface{}{"msg": "sandbox-smoke-test"}
			result, err := exec.Execute(ctx, "echo", params)
			if err != nil {
				return fmt.Errorf("smoke test: %w", err)
			}

			fmt.Printf("Smoke test passed: %v\n", result)
			return nil
		},
	}
}

func newSandboxCleanupCmd(bootLoader func() (*bootstrap.Result, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "cleanup",
		Short: "Remove orphaned sandbox containers",
		Long:  "Find and remove Docker containers with label lango.sandbox=true.",
		RunE: func(cmd *cobra.Command, args []string) error {
			dr, err := sandbox.NewDockerRuntime()
			if err != nil {
				return fmt.Errorf("docker unavailable: %w", err)
			}

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			if !dr.IsAvailable(ctx) {
				return fmt.Errorf("docker daemon is not reachable")
			}

			if err := dr.Cleanup(ctx, ""); err != nil {
				return fmt.Errorf("cleanup: %w", err)
			}

			fmt.Println("Orphaned sandbox containers cleaned up.")
			return nil
		},
	}
}
