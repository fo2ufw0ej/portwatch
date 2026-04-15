package scanner

// Diff holds the result of comparing two port snapshots.
type Diff struct {
	// Opened contains ports present in current but not in previous.
	Opened []int
	// Closed contains ports present in previous but not in current.
	Closed []int
}

// ComputeDiff compares two slices of open ports and returns what changed.
// Both slices are expected to contain unique, sorted port numbers.
func ComputeDiff(previous, current []int) Diff {
	prevSet := toSet(previous)
	currSet := toSet(current)

	var diff Diff

	for port := range currSet {
		if !prevSet[port] {
			diff.Opened = append(diff.Opened, port)
		}
	}

	for port := range prevSet {
		if !currSet[port] {
			diff.Closed = append(diff.Closed, port)
		}
	}

	return diff
}

func toSet(ports []int) map[int]bool {
	s := make(map[int]bool, len(ports))
	for _, p := range ports {
		s[p] = true
	}
	return s
}
