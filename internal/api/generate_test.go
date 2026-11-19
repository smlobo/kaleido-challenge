package api

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// api.gen.go is committed, so the service builds without the generator.
// This test regenerates it from openapi.yaml, with the generator version
// pinned in go.mod, and fails if the result differs: someone changed the spec
// and forgot to run go generate.
func TestGeneratedCodeMatchesSpec(t *testing.T) {
	if testing.Short() {
		t.Skip("runs the code generator")
	}
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go not on PATH")
	}
	// The config's output: setting wins over -o, so point a copy of it at a
	// temporary file instead.
	dir := t.TempDir()
	out := filepath.Join(dir, "api.gen.go")
	cfg, err := os.ReadFile("oapi-codegen.yaml")
	if err != nil {
		t.Fatal(err)
	}
	tmpCfg := filepath.Join(dir, "oapi-codegen.yaml")
	cfg = bytes.Replace(cfg, []byte("output: api.gen.go"), []byte("output: "+out), 1)
	if err := os.WriteFile(tmpCfg, cfg, 0o600); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(goBin, "tool", "oapi-codegen", "-config", tmpCfg, "openapi.yaml")
	if msg, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("oapi-codegen: %v\n%s", err, msg)
	}
	want, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile("api.gen.go")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatal("api.gen.go does not match openapi.yaml; run: go generate ./internal/api")
	}
}
