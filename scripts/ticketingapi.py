"""Minimal HTTP helpers for the ticketing service, shared by the scripts in
this folder. Standard library only."""

import json
import os
import sys
import time
import urllib.error
import urllib.request

DEFAULT_URL = os.environ.get("TICKETING_URL", "http://localhost:8080")


class APIError(RuntimeError):
    def __init__(self, message: str, status: int | None = None):
        super().__init__(message)
        self.status = status  # the HTTP status, or None if the service was not reached

    @property
    def retryable(self) -> bool:
        """True if sending the same request again might succeed: the service
        was unreachable or had an internal error, not a refusal."""
        return self.status is None or self.status >= 500


def add_token_arg(parser, env_var: str, who: str) -> None:
    """Add --token, which falls back to the environment variable env_var."""
    parser.add_argument("--token", metavar="TOKEN",
                        help=f"{who}'s token (default: ${env_var}). A token given here lands in "
                             "your shell history and the process list; the variable does not")


def resolve_token(arg: str | None, env_var: str) -> str:
    """The --token value if given, else the environment variable, else ''."""
    return (arg or os.environ.get(env_var, "")).strip()


def request(method: str, base_url: str, path: str, body: dict | None = None,
            token: str | None = None, headers: dict | None = None) -> dict:
    """Send a request and return the JSON answer, or raise APIError with the
    service's reason. token, if given, is sent as a bearer token."""
    h = {"Content-Type": "application/json"} if body is not None else {}
    if token:
        h["Authorization"] = "Bearer " + token
    h.update(headers or {})
    data = json.dumps(body).encode() if body is not None else None
    req = urllib.request.Request(base_url.rstrip("/") + path, data, h, method=method)
    try:
        with urllib.request.urlopen(req, timeout=30) as resp:
            return json.load(resp)
    except urllib.error.HTTPError as e:
        # The service answers errors as JSON problems: {"title", "status", "code", "detail"}.
        try:
            problem = json.load(e)
            reason = f"{problem.get('title', e.reason)}: {problem.get('detail', '')}".rstrip(": ")
        except ValueError:
            reason = e.reason
        raise APIError(f"{method} {path} failed with {e.code}: {reason}", e.code) from None
    except urllib.error.URLError as e:
        raise APIError(f"cannot reach {base_url} ({e.reason}); is the service running?") from None


def post_json(base_url: str, path: str, body: dict, token: str | None = None,
              headers: dict | None = None) -> dict:
    return request("POST", base_url, path, body, token, headers)


def get_json(base_url: str, path: str, token: str | None = None) -> dict:
    return request("GET", base_url, path, token=token)


FINAL_STATUSES = {"mined", "reverted", "failed"}


def wait_for_operation(base_url: str, op_id: str, token: str, seconds: float) -> dict:
    """Poll GET /operations/{id} until the transaction is final (mined,
    reverted or failed) or seconds pass, and return the last answer."""
    deadline = time.monotonic() + seconds
    op = get_json(base_url, "/operations/" + op_id, token=token)
    while op["status"] not in FINAL_STATUSES and time.monotonic() < deadline:
        time.sleep(1)
        op = get_json(base_url, "/operations/" + op_id, token=token)
    return op


def print_outcome(op: dict | None, seconds: float) -> int:
    """Print how an operation ended and return the script's exit status:
    0 mined (or not waited for), 1 otherwise."""
    if op is None:
        return 0
    line = f"tx          {op['status']}"
    if op.get("block_number"):
        line += f" in block {op['block_number']}"
    print(line)
    if op.get("revert_reason"):
        print(f"reverted    {op['revert_reason']}")
    if op.get("last_error"):
        print(f"last error  {op['last_error']}")
    sys.stdout.flush()  # keep the lines above ahead of any message on stderr
    if op["status"] == "mined":
        return 0
    if op["status"] in ("reverted", "failed"):
        print("The transaction did not go through; the ticket is as it was.", file=sys.stderr)
    else:
        print(f"Still {op['status']} after {seconds:g} s; is Besu producing blocks?", file=sys.stderr)
    return 1
