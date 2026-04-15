package reporter

import "fmt"

// PortListString formats a slice of port numbers as a compact string.
// e.g. [80 443 8080] → "80, 443, 8080" or "(none)" for an empty slice.
func PortListString(ports []int) string {
	if len(ports) == 0 {
		return "(none)"
	}
	s := ""
	for i, p := range ports {
		if i > 0 {
			s += ", "
		}
		s += fmt.Sprintf("%d", p)
	}
	return s
}

// SummaryLine builds a single-line summary string without writing to any
// writer, useful for logging or embedding in alert messages.
func SummaryLine(opened, closed []int) string {
	if len(opened) == 0 && len(closed) == 0 {
		return "no port changes detected"
	}
	return fmt.Sprintf("opened: [%s]  closed: [%s]",
		PortListString(opened), PortListString(closed))
}
