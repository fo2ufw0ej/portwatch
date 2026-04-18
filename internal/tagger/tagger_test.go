package tagger_test

import (
	"testing"

	"github.com/patrickward/portwatch/internal/tagger"
)

func TestLookup_WellKnown(t *testing.T) {
	tr := tagger.New(nil)
	cases := []struct {
		port int
		want string
	}{
		{22, "ssh"},
		{80, "http"},
		{443, "https"},
		{3306, "mysql"},
		{5432, "postgres"},
	}
	for _, tc := range cases {
		got := tr.Lookup(tc.port)
		if got != tc.want {
			t.Errorf("Lookup(%d) = %q; want %q", tc.port, got, tc.want)
		}
	}
}

func TestLookup_Unknown_FallsBack(t *testing.T) {
	tr := tagger.New(nil)
	got := tr.Lookup(9999)
	if got != "port/9999" {
		t.Errorf("expected fallback, got %q", got)
	}
}

func TestLookup_CustomOverridesWellKnown(t *testing.T) {
	tr := tagger.New(map[int]string{80: "my-app"})
	got := tr.Lookup(80)
	if got != "my-app" {
		t.Errorf("expected custom tag, got %q", got)
	}
}

func TestLookup_CustomNewPort(t *testing.T) {
	tr := tagger.New(map[int]string{12345: "internal-api"})
	got := tr.Lookup(12345)
	if got != "internal-api" {
		t.Errorf("expected custom tag, got %q", got)
	}
}

func TestTagAll_ReturnsTags(t *testing.T) {
	tr := tagger.New(nil)
	ports := []int{22, 80, 9999}
	tags := tr.TagAll(ports)
	if len(tags) != 3 {
		t.Fatalf("expected 3 tags, got %d", len(tags))
	}
	if tags[0].Service != "ssh" {
		t.Errorf("expected ssh, got %q", tags[0].Service)
	}
	if tags[1].Service != "http" {
		t.Errorf("expected http, got %q", tags[1].Service)
	}
	if tags[2].Service != "port/9999" {
		t.Errorf("expected port/9999, got %q", tags[2].Service)
	}
}

func TestTagAll_Empty(t *testing.T) {
	tr := tagger.New(nil)
	tags := tr.TagAll([]int{})
	if len(tags) != 0 {
		t.Errorf("expected empty slice, got %d elements", len(tags))
	}
}
