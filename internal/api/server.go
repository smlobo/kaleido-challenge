// Package api is the HTTP layer: openapi.yaml is the contract, api.gen.go is
// generated from it (go generate ./internal/api), and this file turns
// requests into calls on internal/ticketing and its answers into JSON.
package api

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/smlobo/kaleido-challenge/internal/ticketing"
)

//go:embed openapi.yaml
var spec []byte

// Server implements the generated ServerInterface.
type Server struct {
	svc *ticketing.Service
	log *slog.Logger
}

var _ ServerInterface = (*Server)(nil)

// Register adds the API, its spec at /openapi.yaml and Swagger UI at /docs
// to mux.
func Register(mux *http.ServeMux, svc *ticketing.Service, log *slog.Logger) {
	if log == nil {
		log = slog.Default()
	}
	s := &Server{svc: svc, log: log}
	HandlerWithOptions(s, StdHTTPServerOptions{
		BaseRouter:  mux,
		Middlewares: []MiddlewareFunc{s.authenticate},
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			writeProblem(w, http.StatusBadRequest, "invalid_request", "Invalid request", err.Error())
		},
	})
	mux.HandleFunc("GET /openapi.yaml", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		_, _ = w.Write(spec)
	})
	mux.HandleFunc("GET /docs", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = io.WriteString(w, docsPage)
	})
}

// Swagger UI from a pinned CDN version, so the binary carries no UI assets.
const docsPage = `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Ticketing API</title>
<link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui.css">
</head>
<body>
<div id="ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui-bundle.js"></script>
<script>SwaggerUIBundle({ url: "/openapi.yaml", dom_id: "#ui" });</script>
</body>
</html>
`

type principalKey struct{}

// authenticate resolves the bearer token, if any, into the request context.
// A token that is present but unknown is refused here; no token means an
// anonymous caller, and each handler decides whether that is enough.
func (s *Server) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string
		if h := r.Header.Get("Authorization"); h != "" {
			scheme, t, ok := strings.Cut(h, " ")
			if !ok || !strings.EqualFold(scheme, "Bearer") || strings.TrimSpace(t) == "" {
				writeProblem(w, http.StatusUnauthorized, "unauthorized", "Unauthorized", "Use Authorization: Bearer <token>.")
				return
			}
			token = strings.TrimSpace(t)
		}
		p, err := s.svc.Authenticate(r.Context(), token)
		if err != nil {
			s.fail(w, r, err)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), principalKey{}, p)))
	})
}

func caller(r *http.Request) ticketing.Principal {
	p, _ := r.Context().Value(principalKey{}).(ticketing.Principal)
	return p
}

// fail writes err as a problem. Refusals keep their code; anything else is a
// 500 whose details go to the log, not the client.
func (s *Server) fail(w http.ResponseWriter, r *http.Request, err error) {
	var e *ticketing.Error
	if !errors.As(err, &e) {
		if r.Context().Err() == nil {
			s.log.Error("request failed", "method", r.Method, "path", r.URL.Path, "err", err)
		}
		writeProblem(w, http.StatusInternalServerError, "internal", "Internal error", "Something went wrong; it has been logged.")
		return
	}
	status := map[ticketing.Kind]int{
		ticketing.Invalid:      http.StatusBadRequest,
		ticketing.Unauthorized: http.StatusUnauthorized,
		ticketing.Forbidden:    http.StatusForbidden,
		ticketing.NotFound:     http.StatusNotFound,
		ticketing.Conflict:     http.StatusConflict,
		ticketing.Rejected:     http.StatusUnprocessableEntity,
	}[e.Kind]
	writeProblem(w, status, e.Code, e.Title, e.Detail)
}

func writeProblem(w http.ResponseWriter, status int, code, title, detail string) {
	p := Problem{Status: status, Code: code, Title: title}
	if detail != "" {
		p.Detail = &detail
	}
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(p)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// decode reads a JSON body of at most 64 KiB, refusing unknown fields.
func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	d := json.NewDecoder(http.MaxBytesReader(w, r.Body, 64<<10))
	d.DisallowUnknownFields()
	if err := d.Decode(v); err != nil {
		writeProblem(w, http.StatusBadRequest, "invalid_json", "Invalid JSON", err.Error())
		return false
	}
	return true
}

// ---- Accounts ----

func (s *Server) CreateOrganizer(w http.ResponseWriter, r *http.Request) {
	var in NewOrganizer
	if !decode(w, r, &in) {
		return
	}
	o, err := s.svc.CreateOrganizer(r.Context(), caller(r), in.Name)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, CreatedOrganizer{
		Id: uuid.MustParse(o.ID), Name: o.Name, Token: o.Token,
		BoxOfficeAddress: o.BoxOffice, GateAddress: o.Gate, PayoutAddress: o.Payout,
	})
}

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	var in NewUser
	if !decode(w, r, &in) {
		return
	}
	u, err := s.svc.CreateUser(r.Context(), in.Subject)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, CreatedUser{Id: uuid.MustParse(u.ID), Subject: u.Subject, WalletAddress: u.Wallet, Token: u.Token})
}

// ---- Events ----

func (s *Server) CreateEvent(w http.ResponseWriter, r *http.Request, params CreateEventParams) {
	var in NewEvent
	if !decode(w, r, &in) {
		return
	}
	e, err := s.svc.CreateEvent(r.Context(), caller(r), ticketing.NewEvent{
		Name: in.Name, Symbol: in.Symbol, StartsAt: in.StartsAt, MaxSupply: in.MaxSupply,
		FaceValue: in.FaceValue, MaxResaleBps: in.MaxResaleBps,
	}, params.IdempotencyKey)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusAccepted, eventJSON(e))
}

func (s *Server) GetEvent(w http.ResponseWriter, r *http.Request, id EventID) {
	e, err := s.svc.GetEvent(r.Context(), id.String())
	if err != nil {
		s.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, eventJSON(e))
}

func eventJSON(e ticketing.Event) Event {
	return Event{
		Id: uuid.MustParse(e.ID), OrganizerId: uuid.MustParse(e.OrganizerID), Name: e.Name, Symbol: e.Symbol,
		StartsAt: e.StartsAt.UTC(), MaxSupply: e.MaxSupply, FaceValue: e.FaceValue, MaxResaleBps: e.MaxResaleBps,
		Status: EventStatus(e.Status), ContractAddress: e.ContractAddress,
		TicketsIssued: e.TicketsIssued, TicketsSold: e.TicketsSold, DeployOperationId: uuid.MustParse(e.DeployOperationID),
	}
}

// ---- Tickets ----

func (s *Server) BuyTicket(w http.ResponseWriter, r *http.Request, id EventID, params BuyTicketParams) {
	t, err := s.svc.BuyTicket(r.Context(), caller(r), id.String(), params.IdempotencyKey)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	s.writeTicket(w, r, http.StatusAccepted, t)
}

func (s *Server) GetTicket(w http.ResponseWriter, r *http.Request, id TicketID) {
	t, err := s.svc.GetTicket(r.Context(), id.String())
	if err != nil {
		s.fail(w, r, err)
		return
	}
	s.writeTicket(w, r, http.StatusOK, t)
}

func (s *Server) CreatePass(w http.ResponseWriter, r *http.Request, id TicketID) {
	p, err := s.svc.CreatePass(r.Context(), caller(r), id.String())
	if err != nil {
		s.fail(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, CheckInPass{
		TicketId: uuid.MustParse(p.TicketID), TokenId: p.TokenID, ContractAddress: p.Contract,
		ChainId: p.ChainID.Int64(), Deadline: p.Deadline.Unix(), ExpiresAt: p.Deadline.UTC(),
		Signature: p.Signature,
	})
}

func (s *Server) CheckIn(w http.ResponseWriter, r *http.Request, id TicketID, params CheckInParams) {
	var in CheckInRequest
	if !decode(w, r, &in) {
		return
	}
	pass := ticketing.PassInput{Deadline: in.Deadline, Signature: in.Signature}
	if in.TicketId != nil {
		pass.TicketID = in.TicketId.String()
	}
	t, err := s.svc.CheckIn(r.Context(), caller(r), id.String(), pass, params.IdempotencyKey)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	s.writeTicket(w, r, http.StatusAccepted, t)
}

func (s *Server) writeTicket(w http.ResponseWriter, r *http.Request, status int, t ticketing.Ticket) {
	var record map[string]any
	if err := json.Unmarshal(t.Record, &record); err != nil {
		s.fail(w, r, err)
		return
	}
	out := Ticket{
		Id: uuid.MustParse(t.ID), EventId: uuid.MustParse(t.EventID), TokenId: t.TokenID,
		Status: TicketStatus(t.Status), OwnerAddress: t.OwnerAddress, Record: record, RecordHash: t.RecordHash,
		MetadataUrl: t.MetadataURL, MintOperationId: uuid.MustParse(t.MintOperationID),
	}
	if t.CheckInOperationID != nil {
		id := uuid.MustParse(*t.CheckInOperationID)
		out.CheckinOperationId = &id
	}
	writeJSON(w, status, out)
}

func (s *Server) GetMetadata(w http.ResponseWriter, r *http.Request, eventID openapi_types.UUID, tokenID string) {
	m, err := s.svc.GetMetadata(r.Context(), eventID.String(), tokenID)
	if err != nil {
		s.fail(w, r, err)
		return
	}
	var record map[string]any
	if err := json.Unmarshal(m.Record, &record); err != nil {
		s.fail(w, r, err)
		return
	}
	out := Metadata{Name: m.Name, Description: m.Description, Record: record, RecordHash: m.RecordHash}
	for _, a := range m.Attributes {
		out.Attributes = append(out.Attributes, struct {
			TraitType string `json:"trait_type"`
			Value     string `json:"value"`
		}{a[0], a[1]})
	}
	writeJSON(w, http.StatusOK, out)
}

// ---- Operations ----

func (s *Server) GetOperation(w http.ResponseWriter, r *http.Request, id openapi_types.UUID) {
	op, err := s.svc.GetOperation(r.Context(), caller(r), id.String())
	if err != nil {
		s.fail(w, r, err)
		return
	}
	out := Operation{
		Id: id, Status: OperationStatus(op.Status), FromAddress: strings.ToLower(op.From.Hex()),
		Nonce: int64(op.Nonce), TxHash: strings.ToLower(op.TxHash.Hex()), Attempts: op.Attempts,
	}
	if op.BlockNumber != nil {
		b := int64(*op.BlockNumber)
		out.BlockNumber = &b
	}
	if op.ContractAddress != nil {
		a := strings.ToLower(op.ContractAddress.Hex())
		out.ContractAddress = &a
	}
	if op.Revert != "" {
		out.RevertReason = &op.Revert
	}
	if op.LastError != "" {
		out.LastError = &op.LastError
	}
	writeJSON(w, http.StatusOK, out)
}
