package fence_test

import (
	"testing"

	"github.com/guitarkeegan/portwatch/internal/fence"
)

func TestDefaultConfig_IsPermissive(t *testing.T) {
	cfg := fence.DefaultConfig()
	g, err := fence.New(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !g.Allow(80) || !g.Allow(443) || !g.Allow(9999) {
		t.Fatal("default guard should allow all ports")
	}
}

func TestAllow_Denylist(t *testing.T) {
	g, _ := fence.New(fence.Config{Denylist: []int{22, 23}})
	if g.Allow(22) {
		t.Error("port 22 should be denied")
	}
	if !g.Allow(80) {
		t.Error("port 80 should be allowed")
	}
}

func TestAllow_Allowlist(t *testing.T) {
	g, _ := fence.New(fence.Config{Allowlist: []int{80, 443}})
	if !g.Allow(80) {
		t.Error("port 80 should be allowed")
	}
	if g.Allow(8080) {
		t.Error("port 8080 should not be allowed")
	}
}

func TestAllow_DenylistOverridesAllowlist(t *testing.T) {
	g, _ := fence.New(fence.Config{
		Allowlist: []int{80, 443},
		Denylist:  []int{80},
	})
	if g.Allow(80) {
		t.Error("denylist should override allowlist for port 80")
	}
	if !g.Allow(443) {
		t.Error("port 443 should still be allowed")
	}
}

func TestFilter_ReturnsSubset(t *testing.T) {
	g, _ := fence.New(fence.Config{Allowlist: []int{80, 443, 8080}})
	got := g.Filter([]int{22, 80, 443, 9090})
	if len(got) != 2 {
		t.Fatalf("expected 2 ports, got %d: %v", len(got), got)
	}
}

func TestValidate_InvalidPort(t *testing.T) {
	cfg := fence.Config{Allowlist: []int{0}}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for port 0")
	}
}

func TestNew_InvalidConfig_ReturnsError(t *testing.T) {
	_, err := fence.New(fence.Config{Denylist: []int{99999}})
	if err == nil {
		t.Fatal("expected error for out-of-range port")
	}
}
