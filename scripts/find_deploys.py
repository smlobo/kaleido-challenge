#!/usr/bin/env python3
"""List contract-creation transactions for a compiled contract on the dev chain.

A creation has no "to"; it matches the contract if its input starts with the
contract's creation bytecode (constructor args follow the bytecode).

Usage: scripts/find_deploys.py [contract-name] [--rpc URL] [--batch N]
Standard library only.
"""

import argparse
import json
import sys
import urllib.request
from pathlib import Path

BUILD_DIR = Path(__file__).resolve().parent.parent / "contracts" / "build"


class RPC:
    def __init__(self, url: str):
        self.url = url
        self.next_id = 0

    def batch(self, calls: list[tuple[str, list]]) -> list:
        """Send several JSON-RPC calls in one HTTP request; results in call order."""
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
                raise RuntimeError(f"{r['method']}: {reply['error']}")
            results.append(reply["result"])
        return results

    def call(self, method: str, params: list):
        return self.batch([(method, params)])[0]


SELECTOR_NAME = "0x06fdde03"    # name()
SELECTOR_SYMBOL = "0x95d89b41"  # symbol()


def decode_string(result: str) -> str:
    """Decode an ABI-encoded string return value (offset, length, bytes)."""
    data = bytes.fromhex(result[2:])
    offset = int.from_bytes(data[0:32], "big")
    length = int.from_bytes(data[offset:offset + 32], "big")
    return data[offset + 32:offset + 32 + length].decode("utf-8", "replace")


def read_metadata(rpc: RPC, address: str) -> tuple[str, str]:
    """Call name() and symbol() (ERC-20 / ERC-721 metadata); '?' if unsupported."""
    out = []
    for selector in (SELECTOR_NAME, SELECTOR_SYMBOL):
        try:
            out.append(decode_string(rpc.call("eth_call", [{"to": address, "data": selector}, "latest"])))
        except (RuntimeError, ValueError, IndexError):
            out.append("?")
    return out[0], out[1]


def main() -> int:
    p = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    p.add_argument("name", nargs="?", default="DemoNFT", help="contract name in contracts/build")
    p.add_argument("--rpc", default="http://localhost:8545")
    p.add_argument("--batch", type=int, default=200, help="blocks fetched per HTTP request")
    args = p.parse_args()

    bin_path = BUILD_DIR / f"{args.name}.bin"
    if not bin_path.exists():
        print(f"no bytecode at {bin_path}; run contracts/build.sh", file=sys.stderr)
        return 1
    bytecode = "0x" + bin_path.read_text().strip().lower()

    rpc = RPC(args.rpc)
    latest = int(rpc.call("eth_blockNumber", []), 16)
    print(f"scanning blocks 0..{latest} for {args.name} creations", file=sys.stderr)

    found = 0
    for start in range(0, latest + 1, args.batch):
        numbers = range(start, min(start + args.batch, latest + 1))
        blocks = rpc.batch([("eth_getBlockByNumber", [hex(n), True]) for n in numbers])
        creations = [
            (int(block["number"], 16), tx)
            for block in blocks
            for tx in block["transactions"]
            if tx.get("to") is None and tx["input"].lower().startswith(bytecode)
        ]
        if not creations:
            continue
        receipts = rpc.batch([("eth_getTransactionReceipt", [tx["hash"]]) for _, tx in creations])
        for (number, tx), receipt in zip(creations, receipts):
            status = "ok" if receipt["status"] == "0x1" else "reverted"
            label = ""
            if status == "ok":
                name, symbol = read_metadata(rpc, receipt["contractAddress"])
                label = f"  {name!r} ({symbol})"
            print(f"block {number}  tx {tx['hash']}  from {tx['from']}  "
                  f"contract {receipt['contractAddress']}  {status}{label}")
            found += 1

    print(f"{found} found", file=sys.stderr)
    return 0


if __name__ == "__main__":
    sys.exit(main())
