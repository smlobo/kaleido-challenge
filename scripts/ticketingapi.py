"""Minimal HTTP helpers for the ticketing service, shared by the scripts in
this folder. Standard library only."""

import json
import os
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
