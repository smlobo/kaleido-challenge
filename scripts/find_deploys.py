#!/usr/bin/env python3
"""List contract-creation transactions for a compiled contract on the dev chain.

A creation has no "to"; it matches the contract if its input starts with the
contract's creation bytecode (constructor args follow the bytecode).

Usage: scripts/find_deploys.py [contract-name] [--rpc URL] [--batch N]
Standard library only.
"""

import argparse
import sys
from pathlib import Path

from ethrpc import DEFAULT_RPC, RPC, read_metadata

BUILD_DIR = Path(__file__).resolve().parent.parent / "contracts" / "build"


def main() -> int:
    p = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    p.add_argument("name", nargs="?", default="DemoNFT", help="contract name in contracts/build")
    p.add_argument("--rpc", default=DEFAULT_RPC)
    p.add_argument("--batch", type=int, default=200, help="blocks fetched per HTTP request")
    args = p.parse_args()

    bin_path = BUILD_DIR / f"{args.name}.bin"
    if not bin_path.exists():
        print(f"no bytecode at {bin_path}; run contracts/build.sh", file=sys.stderr)
        return 1
    bytecode = "0x" + bin_path.read_text().strip().lower()

    rpc = RPC(args.rpc)
    latest = rpc.block_number()
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
