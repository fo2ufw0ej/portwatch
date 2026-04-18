// Package labelmap maps port numbers to human-readable service labels.
package labelmap

import "fmt"

// LabelMap holds port-to-label mappings.
type LabelMap struct {
	custom   map[int]string
	builtin  map[int]string
}

// New returns a LabelMap seeded with well-known port labels.
func New(custom map[int]string) *LabelMap {
	builtin := map[int]string{
		21:   "ftp",
		22:   "ssh",
		23:   "telnet",
		25:   "smtp",
		53:   "dns",
		80:   "http",
		110:  "pop3",
		143:  "imap",
		443:  "https",
		3306: "mysql",
		5432: "postgres",
		6379: "redis",
		8080: "http-alt",
		8443: "https-alt",
		27017: "mongodb",
	}
	if custom == nil {
		custom = map[int]string{}
	}
	return &LabelMap{custom: custom, builtin: builtin}
}

// Lookup returns the label for a port, preferring custom over builtin.
func (l *LabelMap) Lookup(port int) string {
	if label, ok := l.custom[port]; ok {
		return label
	}
	if label, ok := l.builtin[port]; ok {
		return label
	}
	return fmt.Sprintf("port-%d", port)
}

// LabelAll returns a map of port -> label for all given ports.
func (l *LabelMap) LabelAll(ports []int) map[int]string {
	out := make(map[int]string, len(ports))
	for _, p := range ports {
		out[p] = l.Lookup(p)
	}
	return out
}
