package scanner

// Diff holds the changes between two port scans.
type Diff struct {
	Opened []int // ports that became open
	Closed []int // ports that became closed
}

// HasChanges returns true if any ports opened or closed.
func (d Diff) HasChanges() bool {
	return len(d.Opened) > 0 || len(d.Closed) > 0
}

// ComputeDiff compares a previous and current set of open ports and
// returns the ports that were opened or closed between the two snapshots.
func ComputeDiff(previous, current map[int]struct{}) Diff {
	d := Diff{}

	// Find newly opened ports (in current but not in previous).
	for port := range current {
		if _, existed := previous[port]; !existed {
			d.Opened = append(d.Opened, port)
		}
	}

	// Find newly closed ports (in previous but not in current).
	for port := range previous {
		if _, exists := current[port]; !exists {
			d.Closed = append(d.Closed, port)
		}
	}

	return d
}
