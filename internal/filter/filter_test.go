package filter

import (
	"reflect"
	"sort"
	"testing"
)

func sorted(ports []int) []int {
	s := make([]int, len(ports))
	copy(s, ports)
	sort.Ints(s)
	return s
}

func TestNewRule_Empty(t *testing.T) {
	r := NewRule(nil, nil)
	if !r.IsEmpty() {
		t.Error("expected empty rule")
	}
}

func TestApply_NoRestrictions(t *testing.T) {
	r := NewRule(nil, nil)
	input := []int{80, 443, 8080}
	got := r.Apply(input)
	if !reflect.DeepEqual(sorted(got), sorted(input)) {
		t.Errorf("expected %v, got %v", input, got)
	}
}

func TestApply_IgnorePorts(t *testing.T) {
	r := NewRule([]int{22, 8080}, nil)
	input := []int{22, 80, 443, 8080}
	got := sorted(r.Apply(input))
	want := []int{80, 443}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestApply_AllowedPorts(t *testing.T) {
	r := NewRule(nil, []int{80, 443})
	input := []int{22, 80, 443, 8080}
	got := sorted(r.Apply(input))
	want := []int{80, 443}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestApply_IgnoreAndAllowed_IgnoreTakesPrecedence(t *testing.T) {
	// 80 is both allowed and ignored — ignore wins
	r := NewRule([]int{80}, []int{80, 443})
	input := []int{80, 443, 8080}
	got := sorted(r.Apply(input))
	want := []int{443}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestApply_EmptyInput(t *testing.T) {
	r := NewRule([]int{80}, []int{443})
	got := r.Apply([]int{})
	if len(got) != 0 {
		t.Errorf("expected empty result, got %v", got)
	}
}

func TestIsEmpty_WithIgnore(t *testing.T) {
	r := NewRule([]int{22}, nil)
	if r.IsEmpty() {
		t.Error("expected non-empty rule")
	}
}
