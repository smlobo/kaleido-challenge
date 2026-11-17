#!/usr/bin/env python3
"""List ERC-721 mints on the dev chain, with each token's current owner.

A mint is a Transfer event from the zero address. ERC-721 and ERC-20 share the
Transfer signature; ERC-721 indexes tokenId, so its logs have 4 topics and
ERC-20's have 3. Uses eth_getLogs, so it does not read every block.

Usage: scripts/find_minted.py [--contract ADDR ...] [--to ADDR] [--from-block N] [--rpc URL]
Standard library only.
"""

import argparse
import sys

from ethrpc import (DEFAULT_RPC, RPC, RPCError, SELECTOR_OWNER_OF, read_metadata,
                    uint_word, word_to_address)

# keccak256("Transfer(address,address,uint256)")
TRANSFER_TOPIC = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
ZERO_TOPIC = "0x" + "0" * 64

# Besu caps the block range of one eth_getLogs call (--rpc-max-logs-range, default 5000).
LOGS_RANGE = 5000


def address_topic(addr: str) -> str:
    return "0x" + addr.lower().removeprefix("0x").rjust(64, "0")


def main() -> int:
    p = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    p.add_argument("--contract", action="append", default=[],
                   help="only this NFT contract (repeat for several)")
    p.add_argument("--to", help="only tokens minted to this address")
    p.add_argument("--from-block", type=int, default=0)
    p.add_argument("--rpc", default=DEFAULT_RPC)
    args = p.parse_args()

    rpc = RPC(args.rpc)
    latest = rpc.block_number()
    print(f"scanning blocks {args.from_block}..{latest} for ERC-721 mints", file=sys.stderr)

    topics = [TRANSFER_TOPIC, ZERO_TOPIC, address_topic(args.to) if args.to else None]
    logs = []
    for start in range(args.from_block, latest + 1, LOGS_RANGE):
        flt = {"fromBlock": hex(start), "toBlock": hex(min(start + LOGS_RANGE - 1, latest)),
               "topics": topics}
        if args.contract:
            flt["address"] = args.contract
        logs += [log for log in rpc.call("eth_getLogs", [flt]) if len(log["topics"]) == 4]

    # Current owner of each token; ownerOf reverts for a burned token.
    owners = []
    for log in logs:
        data = SELECTOR_OWNER_OF + uint_word(int(log["topics"][3], 16))
        try:
            owners.append(word_to_address(rpc.eth_call(log["address"], data)))
        except RPCError:
            owners.append("burned")

    names: dict[str, str] = {}
    for log, owner in zip(logs, owners):
        contract = log["address"].lower()
        if contract not in names:
            name, symbol = read_metadata(rpc, contract)
            names[contract] = f"{name!r} ({symbol})"
        token_id = int(log["topics"][3], 16)
        to = word_to_address(log["topics"][2])
        moved = "" if owner == to else f"  now {owner}"
        print(f"block {int(log['blockNumber'], 16)}  tx {log['transactionHash']}  "
              f"contract {contract} {names[contract]}  token {token_id}  to {to}{moved}")

    print(f"{len(logs)} found", file=sys.stderr)
    return 0


if __name__ == "__main__":
    sys.exit(main())
