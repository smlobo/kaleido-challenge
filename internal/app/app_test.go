package app

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/smlobo/kaleido-challenge/internal/chain"
	"github.com/smlobo/kaleido-challenge/internal/dbtest"
	"github.com/smlobo/kaleido-challenge/internal/keystore"
)

// These tests start the real service against Postgres (a fresh schema each)
// and the Besu node from docker/docker-compose.yaml, and skip if either is down.

func testConfig(t *testing.T) Config {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	c, err := chain.Dial(ctx, DefaultRPCURL, nil)
	if err != nil {
		t.Skipf("Besu not reachable (%v); run: docker compose -f docker/docker-compose.yaml up -d --wait besu postgres", err)
	}
	c.Close()
	return Config{
		DatabaseURL:     dbtest.SchemaURL(t),
		RPCURL:          DefaultRPCURL,
		ChainID:         big.NewInt(DefaultChainID),
		KEK:             randomKEK(),
		HTTPAddr:        "127.0.0.1:0",
		PublicBaseURL:   "http://localhost:8080",
		StartupTimeout:  5 * time.Second,
		ShutdownTimeout: 5 * time.Second,
	}
}

func randomKEK() []byte {
	k := make([]byte, keystore.KEKSize)
	_, _ = rand.Read(k)
	return k
}

// serve runs a started app in the background; the returned stop shuts it
// down and fails the test if Serve returned an error.
func serve(t *testing.T, a *App) (stop func()) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- a.Serve(ctx) }()
	var once sync.Once
	stop = func() {
		once.Do(func() {
			cancel()
			select {
			case err := <-done:
				if err != nil {
					t.Errorf("Serve: %v", err)
				}
			case <-time.After(10 * time.Second):
				t.Error("Serve did not return after cancel")
			}
		})
	}
	t.Cleanup(stop)
	return stop
}

func getJSON(t *testing.T, url string) (int, map[string]string) {
	t.Helper()
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var m map[string]string
	if err := json.Unmarshal(body, &m); err != nil {
		t.Fatalf("GET %s: %v in %q", url, err, body)
	}
	return resp.StatusCode, m
}

// First start: migrations, a new deployer key, health and readiness, then a
// clean shutdown. Second start on the same database: the same deployer.
func TestStartCreatesDeployerOnceAndServes(t *testing.T) {
	cfg := testConfig(t)

	a, err := Start(context.Background(), cfg, nil)
	if err != nil {
		t.Fatal(err)
	}
	first := a.Deployer()
	if first.Purpose != keystore.Deployer || first.ID == "" {
		t.Fatalf("deployer = %+v", first)
	}
	stop := serve(t, a)
	base := "http://" + a.Addr()
	if code, body := getJSON(t, base+"/healthz"); code != 200 || body["status"] != "ok" {
		t.Fatalf("healthz: %d %v", code, body)
	}
	if code, body := getJSON(t, base+"/readyz"); code != 200 || body["postgres"] != "ok" || body["besu"] != "ok" {
		t.Fatalf("readyz: %d %v", code, body)
	}
	stop()
	if _, err := http.Get(base + "/healthz"); err == nil {
		t.Fatal("still serving after shutdown")
	}

	b, err := Start(context.Background(), cfg, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Close()
	if b.Deployer().ID != first.ID || b.Deployer().Address != first.Address {
		t.Fatalf("second start: deployer %s, want %s", b.Deployer().Address, first.Address)
	}
}

// Instances starting together on a new database share one deployer.
func TestConcurrentFirstStartsShareOneDeployer(t *testing.T) {
	cfg := testConfig(t)
	var wg sync.WaitGroup
	apps := make([]*App, 3)
	errs := make([]error, 3)
	for i := range apps {
		wg.Go(func() {
			c := cfg
			apps[i], errs[i] = Start(context.Background(), c, nil)
		})
	}
	wg.Wait()
	for i, a := range apps {
		if errs[i] != nil {
			t.Fatalf("start %d: %v", i, errs[i])
		}
		defer a.Close()
		if a.Deployer().ID != apps[0].Deployer().ID {
			t.Fatalf("start %d got deployer %s, start 0 got %s", i, a.Deployer().Address, apps[0].Deployer().Address)
		}
	}
}

// A different KEY_ENCRYPTION_KEY than the one the keys were stored under
// stops startup with a clear message, instead of failing at the first deploy.
func TestStartRejectsWrongKEK(t *testing.T) {
	cfg := testConfig(t)
	a, err := Start(context.Background(), cfg, nil)
	if err != nil {
		t.Fatal(err)
	}
	a.Close()

	cfg.KEK = randomKEK()
	_, err = Start(context.Background(), cfg, nil)
	if !errors.Is(err, keystore.ErrDecrypt) || !strings.Contains(err.Error(), "KEY_ENCRYPTION_KEY does not match") {
		t.Fatalf("err = %v", err)
	}
}

// A node on the wrong chain is a configuration error: no retrying.
func TestStartRejectsWrongChainID(t *testing.T) {
	cfg := testConfig(t)
	cfg.ChainID = big.NewInt(1)
	began := time.Now()
	_, err := Start(context.Background(), cfg, nil)
	if err == nil || !strings.Contains(err.Error(), "serves chain 1337, but CHAIN_ID is 1") {
		t.Fatalf("err = %v", err)
	}
	if time.Since(began) > 3*time.Second {
		t.Fatalf("took %s: a wrong chain ID should fail at once", time.Since(began))
	}
}

// With Postgres unreachable, Start retries until StartupTimeout, then names it.
func TestStartGivesUpWhenPostgresIsUnreachable(t *testing.T) {
	cfg := Config{
		DatabaseURL:    "postgres://ticketing:ticketing@127.0.0.1:1/ticketing?sslmode=disable&connect_timeout=1",
		RPCURL:         DefaultRPCURL,
		ChainID:        big.NewInt(DefaultChainID),
		KEK:            randomKEK(),
		HTTPAddr:       "127.0.0.1:0",
		StartupTimeout: 2 * time.Second,
	}
	began := time.Now()
	_, err := Start(context.Background(), cfg, nil)
	if err == nil || !strings.Contains(err.Error(), "connect to Postgres") {
		t.Fatalf("err = %v", err)
	}
	if d := time.Since(began); d < time.Second || d > 6*time.Second {
		t.Fatalf("gave up after %s, want about StartupTimeout (2 s)", d)
	}
}

func TestReadinessReportsFailingCheck(t *testing.T) {
	besuError := "connection refused"
	h := readiness(map[string]func(context.Context) error{
		"postgres": func(context.Context) error { return nil },
		"besu":     func(context.Context) error { return errors.New(besuError) },
	})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/readyz", nil))
	var body map[string]string
	_ = json.Unmarshal(rec.Body.Bytes(), &body)
	if rec.Code != http.StatusServiceUnavailable || body["postgres"] != "ok" || body["besu"] != besuError {
		t.Fatalf("%d %v", rec.Code, body)
	}
}
