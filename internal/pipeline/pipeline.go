// Package pipeline chains scan → filter → diff → alert steps into a single reusable unit.
package pipeline

import (
	"context"
	"fmt"

	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/scanner"
)

// Step is a function that receives the current open ports and returns a
// (possibly modified) list for the next step.
type Step func(ctx context.Context, ports []int) ([]int, error)

// Pipeline executes an ordered chain of Steps.
type Pipeline struct {
	steps []Step
}

// New creates a Pipeline with the provided steps.
func New(steps ...Step) *Pipeline {
	return &Pipeline{steps: steps}
}

// Run executes every step in order, threading the port list through each one.
func (p *Pipeline) Run(ctx context.Context, ports []int) ([]int, error) {
	var err error
	for i, s := range p.steps {
		ports, err = s(ctx, ports)
		if err != nil {
			return nil, fmt.Errorf("pipeline step %d: %w", i, err)
		}
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
	}
	return ports, nil
}

// FilterStep wraps a filter.Rule as a pipeline Step.
func FilterStep(rule *filter.Rule) Step {
	return func(_ context.Context, ports []int) ([]int, error) {
		return rule.Apply(ports), nil
	}
}

// ScanStep returns a Step that ignores the incoming port list and instead
// performs a fresh scan over the supplied range.
func ScanStep(sc *scanner.Scanner) Step {
	return func(ctx context.Context, _ []int) ([]int, error) {
		return sc.OpenPorts()
	}
}
