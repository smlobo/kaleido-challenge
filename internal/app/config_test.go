package app

import (
	"strings"
	"testing"
)

// A valid KEK for config tests only: 32 bytes of 0x01, base64.
const testKEK = "ABCDEFGHIJKLMNOPQRSTUVWXYZEBAQEBAQEBAQEBAQE="

func env(m map[string]string) func(string) string {
	return func(k string) string { return m[k] }
}

func TestFromEnvDefaults(t *testing.T) {
	cfg, err := FromEnv(env(map[string]string{"KEY_ENCRYPTION_KEY": testKEK}))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.DatabaseURL != DefaultDatabaseURL || cfg.RPCURL != DefaultRPCURL ||
		cfg.ChainID.Int64() != DefaultChainID || cfg.HTTPAddr != ":8080" ||
		cfg.PublicBaseURL != "http://localhost:8080" || len(cfg.KEK) != 32 {
		t.Fatalf("defaults: %+v", cfg)
	}
}

func TestFromEnvOverrides(t *testing.T) {
	cfg, err := FromEnv(env(map[string]string{
		"KEY_ENCRYPTION_KEY": testKEK,
		"DATABASE_URL":       "postgres://u:p@db:5432/x",
		"ETH_RPC_URL":        "http://besu:8545",
		"CHAIN_ID":           "2024",
		"HTTP_PORT":          "9000",
		"PUBLIC_BASE_URL":    "https://tickets.example.test/",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.RPCURL != "http://besu:8545" || cfg.ChainID.Int64() != 2024 || cfg.HTTPAddr != ":9000" ||
		cfg.PublicBaseURL != "https://tickets.example.test" {
		t.Fatalf("overrides: %+v", cfg)
	}
}

// Every problem is reported at once, and no message repeats a secret.
func TestFromEnvReportsAllErrors(t *testing.T) {
	_, err := FromEnv(env(map[string]string{
		"KEY_ENCRYPTION_KEY": "c2hvcnQ=", // decodes to 5 bytes
		"DATABASE_URL":       "postgres://user:s3cret@db:notaport/x",
		"ETH_RPC_URL":        "localhost:8545",
		"CHAIN_ID":           "-5",
		"HTTP_PORT":          "70000",
		"PUBLIC_BASE_URL":    "tickets",
	}))
	if err == nil {
		t.Fatal("accepted")
	}
	msg := err.Error()
	for _, want := range []string{"KEY_ENCRYPTION_KEY", "DATABASE_URL", "ETH_RPC_URL", "CHAIN_ID", "HTTP_PORT", "PUBLIC_BASE_URL"} {
		if !strings.Contains(msg, want) {
			t.Errorf("error does not mention %s:\n%s", want, msg)
		}
	}
	for _, secret := range []string{"c2hvcnQ=", "s3cret"} {
		if strings.Contains(msg, secret) {
			t.Errorf("error repeats a secret (%s):\n%s", secret, msg)
		}
	}
}

func TestFromEnvRequiresKEK(t *testing.T) {
	_, err := FromEnv(env(nil))
	if err == nil || !strings.Contains(err.Error(), "KEY_ENCRYPTION_KEY is not set") {
		t.Fatalf("err = %v", err)
	}
}

func TestFromEnvAdminToken(t *testing.T) {
	cfg, err := FromEnv(env(map[string]string{"KEY_ENCRYPTION_KEY": testKEK}))
	if err != nil || cfg.AdminToken != "" {
		t.Fatalf("unset: %q %v", cfg.AdminToken, err)
	}
	cfg, err = FromEnv(env(map[string]string{"KEY_ENCRYPTION_KEY": testKEK, "ADMIN_TOKEN": "0123456789abcdef"}))
	if err != nil || cfg.AdminToken != "0123456789abcdef" {
		t.Fatalf("set: %q %v", cfg.AdminToken, err)
	}
	_, err = FromEnv(env(map[string]string{"KEY_ENCRYPTION_KEY": testKEK, "ADMIN_TOKEN": "XXX"}))
	if err == nil || !strings.Contains(err.Error(), "ADMIN_TOKEN must be at least ") || strings.Contains(err.Error(), "XXX") {
		t.Fatalf("short: %v", err)
	}
}
