package app

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smlobo/kaleido-challenge/internal/chain"
	"github.com/smlobo/kaleido-challenge/internal/contracts/ticketnft"
)

// The milestone 6 thread through the real service, Postgres and Besu:
// organizer → user → event deployed → two tickets minted → one checked in
// with the holder's signed pass,
// checking each step over HTTP and on-chain.
func TestAPIEventTicketCheckInOnBesu(t *testing.T) {
	cfg := testConfig(t)
	cfg.AdminToken = "test-admin-token-0123456789"
	a, err := Start(context.Background(), cfg, nil)
	if err != nil {
		t.Fatal(err)
	}
	serve(t, a)
	c := &client{t: t, base: "http://" + a.Addr()}

	var org struct {
		Token string `json:"token"`
	}
	c.do("POST", "/organizers", cfg.AdminToken, "", `{"name":"Rock Fest Ltd"}`, 201, &org)
	var alice struct {
		Token         string `json:"token"`
		WalletAddress string `json:"wallet_address"`
	}
	c.do("POST", "/users", "", "", `{"subject":"alice"}`, 201, &alice)

	// Missing Idempotency-Key and wrong caller are refused.
	c.do("POST", "/events", org.Token, "", `{}`, 400, nil)
	starts := time.Now().Add(time.Hour).UTC().Format(time.RFC3339)
	body := `{"name":"Rock Fest 2027","symbol":"ROCK27","starts_at":"` + starts +
		`","max_supply":100,"face_value":"100000000","max_resale_bps":11000}`
	c.do("POST", "/events", alice.Token, "e1", body, 403, nil)

	var ev struct {
		ID              string `json:"id"`
		Status          string `json:"status"`
		ContractAddress string `json:"contract_address"`
		DeployOp        string `json:"deploy_operation_id"`
	}
	c.do("POST", "/events", org.Token, "e1", body, 202, &ev)
	// The 202 shows the event as it is when the answer is built: almost
	// always deploying, but the deploy can be mined first. Either is right.
	if ev.Status != "deploying" && ev.Status != "on_sale" {
		t.Fatalf("new event status %q, want deploying or on_sale", ev.Status)
	}
	c.waitFor("/events/"+ev.ID, "on_sale", &ev)

	type ticketJSON struct {
		ID          string  `json:"id"`
		Status      string  `json:"status"`
		TokenID     string  `json:"token_id"`
		Owner       string  `json:"owner_address"`
		RecordHash  string  `json:"record_hash"`
		MetadataURL string  `json:"metadata_url"`
		CheckInOp   *string `json:"checkin_operation_id"`
	}
	var t1, t2 ticketJSON
	c.do("POST", "/events/"+ev.ID+"/tickets", alice.Token, "buy-1", "", 202, &t1)
	c.do("POST", "/events/"+ev.ID+"/tickets", alice.Token, "buy-2", "", 202, &t2)
	c.waitFor("/tickets/"+t1.ID, "active", &t1)
	c.waitFor("/tickets/"+t2.ID, "active", &t2)
	if t1.TokenID != "1" || t2.TokenID != "2" || t1.Owner != strings.ToLower(alice.WalletAddress) ||
		t2.Owner != strings.ToLower(alice.WalletAddress) {
		t.Fatalf("tickets: %+v %+v; alice %s", t1, t2, alice.WalletAddress)
	}

	// On-chain: owner, record hash and token URI match the API.
	nft := ticketnft.NewTicketNFT()
	contract := common.HexToAddress(ev.ContractAddress)
	owner, err := nft.UnpackOwnerOf(call(t, a.chain, contract, nft.PackOwnerOf(big.NewInt(1))))
	if err != nil || !strings.EqualFold(owner.Hex(), alice.WalletAddress) {
		t.Fatalf("ownerOf(1) = %s, %v", owner.Hex(), err)
	}
	info, err := nft.UnpackTicketInfo(call(t, a.chain, contract, nft.PackTicketInfo(big.NewInt(1))))
	if err != nil || common.Hash(info.RecordHash).Hex() != t1.RecordHash || info.FaceValue.String() != "100000000" {
		t.Fatalf("ticketInfo(1) = %+v, %v; API record hash %s", info, err, t1.RecordHash)
	}
	uri, err := nft.UnpackTokenURI(call(t, a.chain, contract, nft.PackTokenURI(big.NewInt(1))))
	if err != nil || uri != t1.MetadataURL {
		t.Fatalf("tokenURI(1) = %q, %v; API %q", uri, err, t1.MetadataURL)
	}
	var meta struct {
		Name       string `json:"name"`
		RecordHash string `json:"record_hash"`
	}
	c.do("GET", strings.TrimPrefix(uri, cfg.PublicBaseURL), "", "", "", 200, &meta)
	if meta.RecordHash != t1.RecordHash || meta.Name != "Rock Fest 2027 #1" {
		t.Fatalf("metadata: %+v", meta)
	}

	// Check in ticket 1. Alice (the holder) makes a pass; the organizer's gate
	// sends it as the check-in body. Burned on-chain, checked_in in the API,
	// once only.
	c.do("POST", "/tickets/"+t1.ID+"/pass", org.Token, "", "", 403, nil)
	var pass1, pass2 json.RawMessage
	c.do("POST", "/tickets/"+t1.ID+"/pass", alice.Token, "", "", 200, &pass1)
	c.do("POST", "/tickets/"+t2.ID+"/pass", alice.Token, "", "", 200, &pass2)
	c.do("POST", "/tickets/"+t1.ID+"/check-in", alice.Token, "in-1", string(pass1), 403, nil)
	c.do("POST", "/tickets/"+t1.ID+"/check-in", org.Token, "in-0", `{}`, 400, nil)
	c.do("POST", "/tickets/"+t1.ID+"/check-in", org.Token, "in-0", string(pass2), 400, nil)
	c.do("POST", "/tickets/"+t1.ID+"/check-in", org.Token, "in-1", string(pass1), 202, &t1)
	c.waitFor("/tickets/"+t1.ID, "checked_in", &t1)
	if _, err := a.chain.Eth().CallContract(context.Background(),
		ethereum.CallMsg{To: &contract, Data: nft.PackOwnerOf(big.NewInt(1))}, nil); err == nil {
		t.Fatal("ownerOf(1) still answers after check-in; the ticket was not burned")
	}
	c.do("POST", "/tickets/"+t1.ID+"/check-in", org.Token, "in-2", string(pass1), 409, nil)

	var op struct {
		Status string `json:"status"`
		Block  int64  `json:"block_number"`
	}
	c.do("GET", "/operations/"+*t1.CheckInOp, org.Token, "", "", 200, &op)
	if op.Status != "mined" || op.Block == 0 {
		t.Fatalf("check-in operation: %+v", op)
	}
	c.do("GET", "/operations/"+*t1.CheckInOp, alice.Token, "", "", 404, nil)

	c.do("GET", "/openapi.yaml", "", "", "", 200, nil)
	c.do("GET", "/docs", "", "", "", 200, nil)
}

func call(t *testing.T, c *chain.Client, to common.Address, data []byte) []byte {
	t.Helper()
	out, err := c.Eth().CallContract(context.Background(), ethereum.CallMsg{To: &to, Data: data}, nil)
	if err != nil {
		t.Fatalf("eth_call: %v", err)
	}
	return out
}

type client struct {
	t    *testing.T
	base string
	last map[string]any
}

// do sends a request, checks the status and decodes the body into out.
func (c *client) do(method, path, token, idemKey, body string, want int, out any) {
	c.t.Helper()
	req, err := http.NewRequest(method, c.base+path, strings.NewReader(body))
	if err != nil {
		c.t.Fatal(err)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	if idemKey != "" {
		req.Header.Set("Idempotency-Key", idemKey)
	}
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.t.Fatal(err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != want {
		c.t.Fatalf("%s %s: %d, want %d: %s", method, path, resp.StatusCode, want, raw)
	}
	c.last = nil
	_ = json.Unmarshal(raw, &c.last)
	if out != nil && bytes.HasPrefix(bytes.TrimSpace(raw), []byte("{")) {
		if err := json.Unmarshal(raw, out); err != nil {
			c.t.Fatalf("%s %s: %v in %s", method, path, err, raw)
		}
	}
}

// waitFor polls a resource until its status is want.
func (c *client) waitFor(path, want string, out any) {
	c.t.Helper()
	deadline := time.Now().Add(30 * time.Second)
	for {
		c.do("GET", path, "", "", "", 200, out)
		if c.last["status"] == want {
			return
		}
		if time.Now().After(deadline) {
			c.t.Fatalf("%s: status %v after 30 s, want %s", path, c.last["status"], want)
		}
		time.Sleep(250 * time.Millisecond)
	}
}
