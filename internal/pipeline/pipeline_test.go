package pipeline_test

import (
	"context"
	"errors"
	"testing"

	"github.com/user/portwatch/internal/pipeline"
)

func addStep(n int) pipeline.Step {
	return func(_ context.Context, ports []int) ([]int, error) {
		return append(ports, n), nil
	}
}

func errorStep(msg string) pipeline.Step {
	return func(_ context.Context, ports []int) ([]int, error) {
		return nil, errors.New(msg)
	}
}

func TestRun_EmptyPipeline(t *testing.T) {
	pl := pipeline.New()
	out, err := pl.Run(context.Background(), []int{1, 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(out))
	}
}

func TestRun_StepsExecuteInOrder(t *testing.T) {
	pl := pipeline.New(addStep(10), addStep(20), addStep(30))
	out, err := pl.Run(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []int{10, 20, 30}
	for i, p := range want {
		if out[i] != p {
			t.Errorf("index %d: want %d got %d", i, p, out[i])
		}
	}
}

func TestRun_StepError_Propagates(t *testing.T) {
	pl := pipeline.New(addStep(1), errorStep("boom"), addStep(2))
	_, err := pl.Run(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "boom" {
		t.Errorf("expected error message %q, got %q", "boom", err.Error())
	}
}

func TestRun_StepError_StopsExecution(t *testing.T) {
	afterErrorCalled := false
	pl := pipeline.New(errorStep("stop here"), func(_ context.Context, ports []int) ([]int, error) {
		afterErrorCalled = true
		return ports, nil
	})
	_, err := pl.Run(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if afterErrorCalled {
		t.Error("step after error should not have been called")
	}
}

func TestRun_CancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	called := false
	pl := pipeline.New(func(c context.Context, ports []int) ([]int, error) {
		called = true
		return ports, nil
	})
	_, err := pl.Run(ctx, nil)
	// first step runs, but context is checked after — err should be ctx.Err()
	if !called {
		t.Fatal("expected step to be called once")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}
