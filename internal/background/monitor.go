package background

import (
	"go.uber.org/zap"
)

// MonitorSummary provides an aggregate view of background task states.
type MonitorSummary struct {
	Total     int `json:"total"`
	Pending   int `json:"pending"`
	Running   int `json:"running"`
	Done      int `json:"done"`
	Failed    int `json:"failed"`
	Cancelled int `json:"cancelled"`
}

// Monitor provides progress tracking for background tasks.
type Monitor struct {
	manager *Manager
	logger  *zap.SugaredLogger
}

// NewMonitor creates a new Monitor that tracks tasks in the given Manager.
func NewMonitor(manager *Manager, logger *zap.SugaredLogger) *Monitor {
	return &Monitor{
		manager: manager,
		logger:  logger,
	}
}

// ActiveCount returns the number of tasks currently pending or running.
func (m *Monitor) ActiveCount() int {
	snapshots := m.manager.List()
	count := 0
	for _, snap := range snapshots {
		if snap.Status == Pending || snap.Status == Running {
			count++
		}
	}
	return count
}

// Summary returns an aggregate summary of all task states.
func (m *Monitor) Summary() MonitorSummary {
	snapshots := m.manager.List()
	var summary MonitorSummary
	summary.Total = len(snapshots)

	for _, snap := range snapshots {
		switch snap.Status {
		case Pending:
			summary.Pending++
		case Running:
			summary.Running++
		case Done:
			summary.Done++
		case Failed:
			summary.Failed++
		case Cancelled:
			summary.Cancelled++
		}
	}

	return summary
}
