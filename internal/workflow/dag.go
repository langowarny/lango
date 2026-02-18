package workflow

import "fmt"

// DAG represents a directed acyclic graph of workflow steps.
type DAG struct {
	steps    map[string]*Step
	children map[string][]string // stepID -> dependent stepIDs
	parents  map[string][]string // stepID -> dependency stepIDs
}

// NewDAG builds a DAG from a slice of workflow steps.
// It returns an error if a circular dependency is detected.
func NewDAG(steps []Step) (*DAG, error) {
	d := &DAG{
		steps:    make(map[string]*Step, len(steps)),
		children: make(map[string][]string, len(steps)),
		parents:  make(map[string][]string, len(steps)),
	}

	for i := range steps {
		s := &steps[i]
		d.steps[s.ID] = s
		d.parents[s.ID] = s.DependsOn
		for _, dep := range s.DependsOn {
			d.children[dep] = append(d.children[dep], s.ID)
		}
	}

	// Verify DAG property via topological sort.
	if _, err := d.TopologicalSort(); err != nil {
		return nil, err
	}

	return d, nil
}

// TopologicalSort returns layers of step IDs that can be executed in parallel.
// Layer 0 contains steps with no dependencies, layer 1 contains steps whose
// dependencies are all in layer 0, and so on.
func (d *DAG) TopologicalSort() ([][]string, error) {
	inDegree := make(map[string]int, len(d.steps))
	for id := range d.steps {
		inDegree[id] = len(d.parents[id])
	}

	var layers [][]string
	remaining := len(d.steps)

	for remaining > 0 {
		var layer []string
		for id, deg := range inDegree {
			if deg == 0 {
				layer = append(layer, id)
			}
		}
		if len(layer) == 0 {
			return nil, fmt.Errorf("circular dependency detected in DAG")
		}

		// Remove processed nodes.
		for _, id := range layer {
			delete(inDegree, id)
			for _, child := range d.children[id] {
				inDegree[child]--
			}
		}

		layers = append(layers, layer)
		remaining -= len(layer)
	}

	return layers, nil
}

// Roots returns step IDs that have no dependencies.
func (d *DAG) Roots() []string {
	var roots []string
	for id := range d.steps {
		if len(d.parents[id]) == 0 {
			roots = append(roots, id)
		}
	}
	return roots
}

// Ready returns step IDs whose dependencies are all in the completed set.
func (d *DAG) Ready(completed map[string]bool) []string {
	var ready []string
	for id := range d.steps {
		if completed[id] {
			continue
		}
		allDone := true
		for _, dep := range d.parents[id] {
			if !completed[dep] {
				allDone = false
				break
			}
		}
		if allDone {
			ready = append(ready, id)
		}
	}
	return ready
}
