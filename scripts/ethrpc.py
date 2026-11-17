"""Minimal JSON-RPC helpers shared by the scripts in this folder. Standard library only."""

import json
import urllib.request

DEFAULT_RPC = "http://localhost:8545"

SELECTOR_NAME = "0x06fdde03"      # name()
SELECTOR_SYMBOL = "0x95d89b41"    # symbol()
SELECTOR_OWNER_OF = "0x6352211e"  # ownerOf(uint256)


class RPCError(RuntimeError):
    pass


class RPC:
    def __init__(self, url: str = DEFAULT_RPC):
        self.url = url
        self.next_id = 0

    def batch(self, calls: list[tuple[str, list]]) -> list:
        """Send several JSON-RPC calls in one HTTP request; results in call order."""
        if not calls:
            return []
        reqs = []
        for method, params in calls:
            self.next_id += 1
            reqs.append({"jsonrpc": "2.0", "id": self.next_id, "method": method, "params": params})
        body = json.dumps(reqs).encode()
        req = urllib.request.Request(self.url, body, {"Content-Type": "application/json"})
        with urllib.request.urlopen(req, timeout=30) as resp:
            replies = {r["id"]: r for r in json.load(resp)}
        results = []
        for r in reqs:
            reply = replies[r["id"]]
            if "error" in reply:
                raise RPCError(f"{r['method']}: {reply['error']}")
            results.append(reply["result"])
        return results

    def call(self, method: str, params: list):
        return self.batch([(method, params)])[0]

    def block_number(self) -> int:
        return int(self.call("eth_blockNumber", []), 16)

    def eth_call(self, to: str, data: str) -> str:
        return self.call("eth_call", [{"to": to, "data": data}, "latest"])


def decode_string(result: str) -> str:
    """Decode an ABI-encoded string return value (offset, length, bytes)."""
    data = bytes.fromhex(result[2:])
    offset = int.from_bytes(data[0:32], "big")
    length = int.from_bytes(data[offset:offset + 32], "big")
    return data[offset + 32:offset + 32 + length].decode("utf-8", "replace")


def word_to_address(word: str) -> str:
    """Take the low 20 bytes of a 32-byte hex word (an indexed address topic or return value)."""
    return "0x" + word[-40:].lower()


HASH_CHARS = 20  # "0x" + 18 hex digits of a 32-byte hash
ADDR_CHARS = 16  # "0x" + 14 hex digits of a 20-byte address


def shorten(value: str, keep: int, full: bool = False) -> str:
    """The first `keep` characters plus '...', or the whole value when full is set."""
    return value if full or len(value) <= keep else value[:keep] + "..."


def uint_word(n: int) -> str:
    """ABI-encode a uint256 as 64 hex chars, no 0x."""
    return f"{n:064x}"


def read_metadata(rpc: RPC, address: str) -> tuple[str, str]:
    """Call name() and symbol() (ERC-20 / ERC-721 metadata); '?' if unsupported."""
    out = []
    for selector in (SELECTOR_NAME, SELECTOR_SYMBOL):
        try:
            out.append(decode_string(rpc.eth_call(address, selector)))
        except (RPCError, ValueError, IndexError):
            out.append("?")
    return out[0], out[1]
