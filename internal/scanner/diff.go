package scanner

import "sort"

// Diff holds the result of comparing two port snapshots.
type Diff struct {
	Opened []int
	Closed []int
}

// HasChanges returns true when at least one port opened or closed.
func (d Diff) HasChanges() bool {
	return len(d.Opened) > 0 || len(d.Closed) > 0
}

// ComputeDiff compares a previous and current list of open ports and
// returns which ports were opened and which were closed.
func ComputeDiff(previous, current []int) Diff {
	prev := toSet(previous)
	curr := toSet(current)

	var opened, closed []int

	for p := range curr {
		if !prev[p] {
			opened = append(opened, p)
		}
	}
	for p := range prev {
		if !curr[p] {
			closed = append(closed, p)
		}
	}

	sort.Ints(opened)
	sort.Ints(closed)

	return Diff{Opened: opened, Closed: closed}
}

func toSet(ports []int) map[int]bool {
	s := make(map[int]bool, len(ports))
	for _, p := range ports {
		s[p] = true
	}
	return s
}
