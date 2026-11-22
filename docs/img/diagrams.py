"""Generates the README diagrams as hand-placed SVGs, light and dark.

Regenerate after editing: python3 docs/img/diagrams.py docs/img
"""
import sys, os
from xml.sax.saxutils import escape

PAL = {
    "light": dict(bg="#ffffff", fg="#1f2328", muted="#59636e", line="#d1d9e0", panel="#f6f8fa",
                  go="#0969da", go_bg="#ddf4ff", pg="#1a7f37", pg_bg="#dafbe1",
                  chain="#8250df", chain_bg="#fbefff", lib="#59636e", lib_bg="#f6f8fa",
                  fe="#bc4c00", fe_bg="#fff1e5", arrow="#59636e"),
    "dark": dict(bg="#0d1117", fg="#e6edf3", muted="#9198a1", line="#3d444d", panel="#151b23",
                 go="#4493f8", go_bg="#0c2d6b", pg="#3fb950", pg_bg="#0f3320",
                 chain="#ab7df8", chain_bg="#2e1a52", lib="#9198a1", lib_bg="#1f242b",
                 fe="#f0883e", fe_bg="#3d2412", arrow="#9198a1"),
}
FONT = "-apple-system,BlinkMacSystemFont,'Segoe UI',Helvetica,Arial,sans-serif"
MONO = "ui-monospace,SFMono-Regular,Menlo,Consolas,monospace"


class SVG:
    def __init__(self, w, h, p):
        self.w, self.h, self.p, self.out = w, h, p, []

    def add(self, s):
        self.out.append(s)

    def rect(self, x, y, w, h, stroke, fill, r=8, dash=False, sw=1.5):
        d = ' stroke-dasharray="6 4"' if dash else ""
        self.add(f'<rect x="{x}" y="{y}" width="{w}" height="{h}" rx="{r}" fill="{fill}" '
                 f'stroke="{stroke}" stroke-width="{sw}"{d}/>')

    def text(self, x, y, s, size=13, color=None, weight=400, anchor="start", mono=False):
        self.add(f'<text x="{x}" y="{y}" font-family="{MONO if mono else FONT}" font-size="{size}" '
                 f'font-weight="{weight}" fill="{color or self.p["fg"]}" text-anchor="{anchor}">{escape(s)}</text>')

    def box(self, x, y, w, h, kind, title, lines=(), mono_last=False, dash=False, center=False):
        p = self.p
        self.rect(x, y, w, h, p[kind], p[kind + "_bg"], dash=dash)
        tx, anc = (x + w / 2, "middle") if center else (x + 12, "start")
        n = 1 + len(lines)
        ty = y + h / 2 - (n - 1) * 8 + 5
        self.text(tx, ty, title, 13.5, weight=600, anchor=anc)
        for i, l in enumerate(lines):
            mono = mono_last and i == len(lines) - 1
            self.text(tx, ty + 17 * (i + 1), l, 11.5 if mono else 12, p["muted"], anchor=anc, mono=mono)

    def panel(self, x, y, w, h, title, sub=None, kind=None, dash=False):
        p = self.p
        self.rect(x, y, w, h, p[kind] if kind else p["line"], p["panel"], r=12, dash=dash)
        self.text(x + 14, y + 24, title, 14, weight=700)
        if sub:
            self.text(x + 14, y + 42, sub, 12, p["muted"])

    def arrow(self, pts, label=None, lpos=None, anchor="middle", both=False):
        p = self.p
        d = "M" + " L".join(f"{a},{b}" for a, b in pts)
        ms = ' marker-start="url(#ah)"' if both else ""
        self.add(f'<path d="{d}" fill="none" stroke="{p["arrow"]}" stroke-width="1.5" '
                 f'marker-end="url(#ah)"{ms}/>')
        if label:
            for i, l in enumerate(label.split("\n")):
                self.text(lpos[0], lpos[1] + 14 * i, l, 11.5, p["muted"], anchor=anchor)

    def svg(self, title):
        p = self.p
        head = (f'<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 {self.w} {self.h}" '
                f'width="{self.w}" height="{self.h}" role="img" aria-label="{escape(title)}">'
                f'<title>{escape(title)}</title>'
                f'<defs><marker id="ah" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="7" '
                f'markerHeight="7" orient="auto-start-reverse"><path d="M0,0 L10,5 L0,10 z" '
                f'fill="{p["arrow"]}"/></marker></defs>'
                f'<rect width="100%" height="100%" fill="{p["bg"]}"/>')
        return head + "".join(self.out) + "</svg>\n"


def legend(s, x, y, items):
    for kind, label in items:
        s.rect(x, y - 10, 14, 14, s.p[kind], s.p[kind + "_bg"], r=3)
        s.text(x + 20, y + 1, label, 12, s.p["muted"])
        x += 28 + len(label) * 6.6


def services(p):
    s = SVG(1000, 600, p)
    # Browser
    s.box(20, 250, 120, 64, "lib", "Browser", ["demo page"], center=True)
    # Frontend
    s.panel(170, 170, 190, 224, "frontend", "nginx", kind="fe")
    s.box(186, 230, 158, 56, "fe", "React page", ["static files"], center=True)
    s.box(186, 310, 158, 56, "fe", "/api proxy", ["→ ticketing"], center=True)
    s.arrow([(140, 270), (163, 270), (163, 258), (184, 258)])
    s.arrow([(140, 296), (163, 296), (163, 338), (184, 338)])
    s.text(152, 412, "HTTP, one origin", 11.5, p["muted"], anchor="middle")

    # Ticketing service, with a second instance behind it
    s.rect(408, 32, 330, 480, p["go"], p["panel"], r=12, dash=True)
    s.panel(398, 22, 330, 480, "ticketing service (Go)", "1..n instances, one binary, env config", kind="go")
    s.box(420, 84, 286, 56, "go", "api", ["OpenAPI handlers, bearer tokens"])
    s.box(420, 166, 286, 56, "go", "ticketing", ["business rules, dry runs, final-status hook"])
    s.box(420, 262, 136, 76, "go", "outbox", ["save, send,", "submitter + receipt"])
    s.box(570, 262, 136, 76, "go", "keystore", ["encrypted keys,", "nonces"])
    s.box(570, 374, 136, 60, "go", "chain", ["JSON-RPC client"])
    s.box(420, 374, 136, 60, "lib", "telemetry", ["OTel spans"])
    s.arrow([(563, 140), (563, 164)])
    s.arrow([(488, 222), (488, 260)])
    s.arrow([(638, 222), (638, 260)])
    s.arrow([(556, 300), (568, 300)])
    s.arrow([(520, 338), (520, 356), (600, 356), (600, 372)])
    s.arrow([(706, 194), (716, 194), (716, 390), (708, 390)])
    s.text(570, 456, "dry run, ticketing → chain:", 11.5, p["muted"])
    s.text(570, 471, "eth_estimateGas", 11.5, p["muted"], mono=True)
    # proxy → api
    s.arrow([(344, 338), (380, 338), (380, 112), (418, 112)])
    s.text(386, 160, "HTTP", 11.5, p["muted"])

    # Postgres (cylinder)
    px, py, pw, ph = 790, 40, 190, 150
    s.add(f'<path d="M{px},{py+14} L{px},{py+ph-14} A{pw/2},14 0 0 0 {px+pw},{py+ph-14} '
          f'L{px+pw},{py+14}" fill="{p["pg_bg"]}" stroke="{p["pg"]}" stroke-width="1.5"/>')
    s.add(f'<ellipse cx="{px+pw/2}" cy="{py+14}" rx="{pw/2}" ry="14" fill="{p["pg_bg"]}" '
          f'stroke="{p["pg"]}" stroke-width="1.5"/>')
    s.text(px + pw / 2, py + 58, "Postgres", 14, weight=700, anchor="middle")
    for i, l in enumerate(["keys, users, organizers", "events, tickets", "transactions (outbox)", "all locking between instances"]):
        s.text(px + pw / 2, py + 80 + 16 * i, l, 12, p["muted"], anchor="middle")
    s.arrow([(728, 115), (788, 115)])
    s.text(758, 104, "SQL", 11.5, p["muted"], anchor="middle")
    s.text(758, 132, "pgx", 11.5, p["muted"], anchor="middle")

    # Besu
    s.panel(770, 236, 210, 210, "Besu", "QBFT, chain 1337, 2 s blocks", kind="chain")
    for i in (2, 1, 0):
        s.rect(790 + i * 8, 290 + i * 8, 160, 60, p["chain"], p["chain_bg"])
    s.text(870, 314, "TicketNFT", 13.5, weight=600, anchor="middle")
    s.text(870, 331, "one contract per event", 12, p["muted"], anchor="middle")
    s.text(784, 392, "roles: box office mints,", 12, p["muted"])
    s.text(784, 408, "gate checks in", 12, p["muted"])
    s.arrow([(706, 418), (750, 418), (750, 330), (768, 330)])
    s.text(744, 466, "eth_sendRawTransaction", 11.5, p["muted"], mono=True)
    s.text(744, 481, "eth_getTransactionReceipt", 11.5, p["muted"], mono=True)

    # Jaeger
    s.box(790, 500, 190, 56, "lib", "Jaeger", ["traces, UI on :16686"], center=True)
    s.arrow([(488, 434), (488, 528), (788, 528)])
    s.text(500, 520, "OTLP spans", 11.5, p["muted"])
    s.arrow([(965, 446), (965, 498)])
    s.text(958, 496, "OTLP", 11.5, p["muted"], anchor="end")

    legend(s, 20, 585, [("go", "Go service"), ("fe", "frontend"), ("pg", "Postgres"),
                         ("chain", "chain"), ("lib", "other")])
    return s.svg("Services: browser, frontend, ticketing, Postgres, Besu, Jaeger")


def concurrency(p):
    s = SVG(1000, 660, p)
    s.panel(20, 20, 580, 600, "one ticketing instance", "a single process", kind="go")
    s.box(40, 76, 540, 58, "go", "main goroutine",
          ["startup (serial) → wait for SIGTERM or an HTTP error → shutdown"])

    # HTTP side
    s.box(40, 168, 260, 50, "go", "HTTP server", ["go a.srv.Serve(ln)"], mono_last=True)
    for i in (2, 1):
        s.rect(40 + i * 6, 246 + i * 6, 248, 96, p["go"], p["go_bg"])
    s.box(40, 246, 248, 96, "go", "request goroutines",
          ["one per request:", "auth, checks, dry run,", "save transaction, fast-path send"])
    s.arrow([(170, 134), (170, 166)])
    s.arrow([(170, 218), (170, 244)])

    # Outbox side
    s.box(320, 168, 260, 50, "go", "outbox.Run", ["go … · own context"])
    s.box(320, 246, 260, 64, "go", "submitter loop (every 1 s)", ["send signed rows,", "resend overdue ones"])
    s.box(320, 322, 260, 64, "go", "receipt loop (every 1 s)", ["record receipts,", "run the final-status hook"])
    s.arrow([(450, 134), (450, 166)])
    s.arrow([(450, 218), (450, 244)])
    s.add(f'<path d="M580,193 L594,193 L594,354 L582,354" fill="none" stroke="{p["arrow"]}" '
          f'stroke-width="1.5" marker-end="url(#ah)"/>')

    s.box(40, 410, 540, 56, "lib", "library goroutines",
          ["OTel batch span exporter · pgxpool health check · signal watcher"], dash=True)

    # Shared clients
    s.text(40, 500, "shared, safe for concurrent use", 12, p["muted"])
    s.box(40, 512, 260, 64, "pg", "pgxpool", ["one connection per", "open transaction"])
    s.box(320, 512, 260, 64, "chain", "chain RPC client", ["10 s timeout per call", "while a row is locked"])

    # Postgres
    s.panel(640, 20, 340, 362, "Postgres", "where the synchronisation lives", kind="pg")
    rows = [
        ("row locks", ["signing_keys: nonces", "events: ticket numbers", "tickets: resale, check-in"]),
        ("work queue", ["transactions … FOR UPDATE SKIP LOCKED"]),
        ("unique constraints", ["idempotency key, (from, nonce),", "one deployer key"]),
        ("advisory lock", ["migrations, one instance at a time"]),
    ]
    y = 76
    for title, lines in rows:
        h = 30 + 17 * len(lines)
        s.box(656, y, 308, h, "pg", title, lines)
        y += h + 12
    # pool → Postgres
    s.arrow([(300, 530), (310, 530), (310, 484), (618, 484), (618, 230), (638, 230)])

    # Other instances
    s.box(640, 420, 340, 58, "go", "other ticketing instances",
          ["the same goroutines, the same Postgres"], dash=True)
    s.arrow([(810, 420), (810, 384)])

    # Besu
    s.box(640, 548, 340, 62, "chain", "Besu", ["orders each key's transactions by nonce"])
    s.arrow([(580, 560), (638, 560)])

    legend(s, 20, 646, [("go", "Go: the service's goroutines"), ("pg", "Postgres"),
                         ("chain", "chain"), ("lib", "started by a library")])
    return s.svg("Concurrency map: goroutines in one instance and the Postgres locks they share")


def write_flow(p):
    s = SVG(1000, 808, p)
    X = {"C": 70, "S": 280, "P": 640, "B": 900}
    heads = [("C", "lib", "Client"), ("S", "go", "ticketing"), ("P", "pg", "Postgres"), ("B", "chain", "Besu")]

    # Phase bands first, so lifelines and arrows draw on top.
    bands = [
        (122, 58, "1  dry run", "no lock held"),
        (192, 112, "2  save", "one Postgres transaction"),
        (316, 112, "3  send (fast path)", "a second transaction"),
        (466, 112, "3b  submitter loop, 1 s", "only if the fast path", "didn't send", True),
        (590, 90, "4  receipt loop, every 1 s", "background goroutine"),
        (690, 104, "5  read the result", "from Postgres only"),
    ]
    for y, h, title, sub, *rest in bands:
        dashed = bool(rest and rest[-1] is True)
        s.rect(10, y, 980, h, p["muted"] if dashed else p["line"], p["panel"], r=8, dash=dashed)
        if rest and isinstance(rest[0], str):
            s.text(84, y + 52, rest[0], 11.5, p["muted"])
        s.text(84, y + 20, title, 12.5, weight=700)
        s.text(84, y + 36, sub, 11.5, p["muted"])

    for k, kind, name in heads:
        s.box(X[k] - 60, 20, 120, 40, kind, name, center=True)
        s.add(f'<line x1="{X[k]}" y1="60" x2="{X[k]}" y2="788" stroke="{p["line"]}" '
              f'stroke-width="1.5" stroke-dasharray="4 4"/>')

    def msg(a, b, y, label, reply=False, mono=False):
        x1, x2 = X[a], X[b]
        x2 += -4 if x2 > x1 else 4
        dash = ' stroke-dasharray="6 4"' if reply else ""
        s.add(f'<path d="M{x1},{y} L{x2},{y}" fill="none" stroke="{p["arrow"]}" stroke-width="1.5" '
              f'marker-end="url(#ah)"{dash}/>')
        s.text((X[a] + X[b]) / 2, y - 7, label, 12, mono=mono, anchor="middle")

    msg("C", "S", 100, "POST, with Idempotency-Key")
    msg("S", "B", 162, "eth_estimateGas: would it revert? then refuse now (4xx)")
    msg("S", "P", 232, "BEGIN · lock the event or ticket row")
    msg("S", "P", 262, "take the next nonce (signing_keys row) · sign in Go")
    msg("S", "P", 292, "INSERT transactions (signed) · COMMIT")
    msg("S", "P", 356, "BEGIN · claim the row FOR UPDATE SKIP LOCKED")
    msg("S", "B", 386, "eth_sendRawTransaction")
    msg("S", "P", 416, "mark submitted · COMMIT")
    msg("S", "C", 448, "202 · operation ID", reply=True)
    msg("S", "P", 500, "claim signed rows, or submitted with no receipt for 10 s")
    msg("S", "B", 530, "eth_sendRawTransaction: oldest nonce first, per key")
    msg("S", "P", 560, "mark submitted · COMMIT")
    msg("S", "B", 630, "eth_getTransactionReceipt")
    msg("S", "P", 664, "mined or reverted · update the ticket or event")
    msg("C", "S", 748, "GET /operations/{id}")
    msg("S", "P", 766, "SELECT (no chain call)")
    msg("S", "C", 784, "status, block number", reply=True)
    return s.svg("The life of a write: dry run, save, send (or the submitter loop), 202, then the receipt loop")


out = sys.argv[1]
os.makedirs(out, exist_ok=True)
for mode, p in PAL.items():
    open(f"{out}/architecture-{mode}.svg", "w").write(services(p))
    open(f"{out}/concurrency-{mode}.svg", "w").write(concurrency(p))
    open(f"{out}/write-flow-{mode}.svg", "w").write(write_flow(p))
print("ok")
