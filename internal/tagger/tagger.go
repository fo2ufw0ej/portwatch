// Package tagger assigns human-readable tags to ports based on well-known
// service mappings (e.g. 22 → "ssh", 443 → "https").
package tagger

import "fmt"

// Tag holds a port number and its resolved service name.
type Tag struct {
	Port    int    `json:"port"`
	Service string `json:"service"`
}

// wellKnown maps port numbers to service names.
var wellKnown = map[int]string{
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

// Tagger resolves port numbers to service tags.
type Tagger struct {
	custom map[int]string
}

// New returns a Tagger optionally extended with custom port→service mappings.
func New(custom map[int]string) *Tagger {
	c := make(map[int]string, len(custom))
	for k, v := range custom {
		c[k] = v
	}
	return &Tagger{custom: c}
}

// Lookup returns the service name for a port, falling back to "port/<n>".
func (t *Tagger) Lookup(port int) string {
	if svc, ok := t.custom[port]; ok {
		return svc
	}
	if svc, ok := wellKnown[port]; ok {
		return svc
	}
	return fmt.Sprintf("port/%d", port)
}

// TagAll returns a Tag slice for the given port list.
func (t *Tagger) TagAll(ports []int) []Tag {
	tags := make([]Tag, len(ports))
	for i, p := range ports {
		tags[i] = Tag{Port: p, Service: t.Lookup(p)}
	}
	return tags
}
