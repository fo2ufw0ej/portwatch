package filter

// Rule defines a port filtering rule.
type Rule struct {
	// IgnorePorts is a set of ports to always exclude from diff alerts.
	IgnorePorts map[int]struct{}
	// AllowedPorts, when non-empty, restricts alerts to only these ports.
	AllowedPorts map[int]struct{}
}

// NewRule constructs a Rule from ignore and allowed port slices.
func NewRule(ignorePorts, allowedPorts []int) *Rule {
	r := &Rule{
		IgnorePorts:  make(map[int]struct{}, len(ignorePorts)),
		AllowedPorts: make(map[int]struct{}, len(allowedPorts)),
	}
	for _, p := range ignorePorts {
		r.IgnorePorts[p] = struct{}{}
	}
	for _, p := range allowedPorts {
		r.AllowedPorts[p] = struct{}{}
	}
	return r
}

// Apply filters a slice of ports according to the rule.
// Ports in IgnorePorts are removed.
// If AllowedPorts is non-empty, only those ports are kept.
func (r *Rule) Apply(ports []int) []int {
	result := make([]int, 0, len(ports))
	for _, p := range ports {
		if _, ignored := r.IgnorePorts[p]; ignored {
			continue
		}
		if len(r.AllowedPorts) > 0 {
			if _, allowed := r.AllowedPorts[p]; !allowed {
				continue
			}
		}
		result = append(result, p)
	}
	return result
}

// IsEmpty returns true when the rule has no restrictions defined.
func (r *Rule) IsEmpty() bool {
	return len(r.IgnorePorts) == 0 && len(r.AllowedPorts) == 0
}
