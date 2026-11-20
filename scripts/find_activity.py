#!/usr/bin/env python3
"""List contract deploys and ERC-721 token activity (mints, listings, resales, burns)
on the dev chain, in chain order, with who signed each transaction.

Deploys: every contract-creation transaction (no "to"). Its type (DemoNFT,
TicketNFT, ...) is the contract in contracts/build whose creation bytecode the
input starts with; "?" if none matches, e.g. a deploy from before a rebuild.

Token activity: ERC-721 Transfer events (4 topics; ERC-20's Transfer has 3).
"from" and "to" are the token's previous and new owner, "by" is the transaction's
signer (left out for a resale, where it is always the buyer, "to"). For TicketNFT, a TicketResold event in the same transaction marks a
transfer as a resale with its price, and a CheckedIn event marks a burn as a
check-in.

Listings: TicketNFT's TicketListed and ListingCancelled events. Listing a ticket
or cancelling a listing moves no token, so there is no Transfer to show; these
lines show the ticket and its holder ("seller"), who is also the signer. --from selects a
seller's listings too; --to leaves listings out, as nothing moves to anyone.

Usage: scripts/find_activity.py [--contract ADDR ...] [--from ADDR] [--to ADDR] [--token ID]
                                [--from-block N] [--full] [--rpc URL]
--from, --to and --token select token activity only, so deploys are left out.
Standard library only.
"""

import argparse
import sys
from pathlib import Path

from ethrpc import (ADDR_CHARS, DEFAULT_RPC, HASH_CHARS, RPC, read_metadata, shorten,
                    uint_word, word_to_address)

BUILD_DIR = Path(__file__).resolve().parent.parent / "contracts" / "build"

# keccak256 of each event signature
TRANSFER_TOPIC = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"  # Transfer(address,address,uint256)
RESOLD_TOPIC = "0xd59a36d8310ffb140214be490baf3fe8961f859a17199ec1bf9015c8b7340011"    # TicketResold(uint256,address,address,uint256)
CHECKED_IN_TOPIC = "0x70790c0ab38e2aca9b5520985fe66db2a3700c3680b61c3853b5554f67a2c31c"  # CheckedIn(uint256,address)
LISTED_TOPIC = "0x22968ec3be3687b180f3544163a034b88d57b48159d9a7e9a1b610fdb2d05f80"      # TicketListed(uint256,address,uint256)
CANCELLED_TOPIC = "0x411aee90354c51b1b04cd563fcab2617142a9d50da19232d888547c8a1b7fd8a"   # ListingCancelled(uint256)
ZERO_TOPIC = "0x" + "0" * 64

# Besu caps the block range of one eth_getLogs call (--rpc-max-logs-range, default 5000).
LOGS_RANGE = 5000
# Blocks fetched per batched HTTP request when scanning for deploys.
BLOCK_BATCH = 200


def address_topic(addr: str) -> str:
    return "0x" + addr.lower().removeprefix("0x").rjust(64, "0")


def load_bytecodes() -> dict[str, str]:
    """Contract name -> lowercase 0x creation bytecode, for every contract in contracts/build."""
    return {p.stem: "0x" + p.read_text().strip().lower() for p in sorted(BUILD_DIR.glob("*.bin"))}


def contract_type(tx_input: str, bytecodes: dict[str, str]) -> str:
    """The built contract whose bytecode the input starts with (longest wins), or '?'."""
    hits = [name for name, code in bytecodes.items() if tx_input.startswith(code)]
    return max(hits, key=lambda n: len(bytecodes[n]), default="?")


def find_deploys(rpc: RPC, first: int, last: int) -> list[dict]:
    """Every contract-creation transaction in [first, last], with its receipt."""
    bytecodes = load_bytecodes()
    creations = []
    for start in range(first, last + 1, BLOCK_BATCH):
        numbers = range(start, min(start + BLOCK_BATCH, last + 1))
        blocks = rpc.batch([("eth_getBlockByNumber", [hex(n), True]) for n in numbers])
        creations += [tx for block in blocks for tx in block["transactions"] if tx.get("to") is None]
    receipts = rpc.batch([("eth_getTransactionReceipt", [tx["hash"]]) for tx in creations])
    return [
        {
            "block": int(tx["blockNumber"], 16),
            "order": (int(tx["transactionIndex"], 16), -1),  # before any log of the same tx
            "tx": tx["hash"],
            "contract": (r.get("contractAddress") or "").lower(),
            "type": contract_type(tx["input"].lower(), bytecodes),
            "signer": tx["from"].lower(),
            "ok": r["status"] == "0x1",
        }
        for tx, r in zip(creations, receipts)
    ]


def get_logs(rpc: RPC, first: int, last: int, topics: list, contracts: list[str]) -> list[dict]:
    """eth_getLogs over [first, last] in LOGS_RANGE chunks."""
    logs = []
    for start in range(first, last + 1, LOGS_RANGE):
        flt = {"fromBlock": hex(start), "toBlock": hex(min(start + LOGS_RANGE - 1, last)), "topics": topics}
        if contracts:
            flt["address"] = contracts
        logs += rpc.call("eth_getLogs", [flt])
    return logs


def main() -> int:
    p = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    p.add_argument("--contract", action="append", default=[],
                   help="only this contract (repeat for several)")
    p.add_argument("--from", dest="sender", metavar="ADDR",
                   help="only tokens moved away from this address, and its listings")
    p.add_argument("--to", help="only tokens moved to this address (leaves listings out)")
    p.add_argument("--token", type=int, help="only this token ID")
    p.add_argument("--from-block", type=int, default=0)
    p.add_argument("--rpc", default=DEFAULT_RPC)
    p.add_argument("--full", action="store_true", help="print full hashes and addresses")
    args = p.parse_args()
    contracts = [c.lower() for c in args.contract]
    token_filter = args.sender or args.to or args.token is not None

    def h(v: str) -> str:
        return shorten(v, HASH_CHARS, args.full)

    def a(v: str) -> str:
        return shorten(v, ADDR_CHARS, args.full)

    rpc = RPC(args.rpc)
    latest = rpc.block_number()
    print(f"scanning blocks {args.from_block}..{latest} for deploys, token activity and listings", file=sys.stderr)

    # Deploys are needed even when filtered out: they give each contract its type.
    deploys = find_deploys(rpc, args.from_block, latest)
    types = {d["contract"]: d["type"] for d in deploys if d["contract"]}

    topics = [
        TRANSFER_TOPIC,
        address_topic(args.sender) if args.sender else None,
        address_topic(args.to) if args.to else None,
        "0x" + uint_word(args.token) if args.token is not None else None,
    ]
    transfers = [log for log in get_logs(rpc, args.from_block, latest, topics, contracts)
                 if len(log["topics"]) == 4]

    # TicketNFT's own events in the same transactions, keyed by (tx, contract, token ID).
    labels: dict[tuple[str, str, int], str] = {}
    if transfers:
        for log in get_logs(rpc, args.from_block, latest, [[RESOLD_TOPIC, CHECKED_IN_TOPIC]], contracts):
            key = (log["transactionHash"], log["address"].lower(), int(log["topics"][1], 16))
            labels[key] = (f"resale price {int(log['data'], 16)}" if log["topics"][0] == RESOLD_TOPIC
                           else "check-in")

    # Listings and cancellations. Both index the token ID; only a listing names
    # the seller, so a cancellation's holder is checked against its signer.
    listings = []
    if not args.to:
        token_topic = "0x" + uint_word(args.token) if args.token is not None else None
        listings = get_logs(rpc, args.from_block, latest, [[LISTED_TOPIC, CANCELLED_TOPIC], token_topic], contracts)

    txs = rpc.batch([("eth_getTransactionByHash", [log["transactionHash"]]) for log in transfers + listings])
    listing_txs = txs[len(transfers):]
    txs = txs[:len(transfers)]

    rows = []  # (block, order, line)
    names: dict[str, str] = {}

    def describe(contract: str) -> str:
        if contract not in names:
            name, symbol = read_metadata(rpc, contract)
            names[contract] = f"{name!r} ({symbol})"
        return names[contract]

    if not token_filter:
        for d in deploys:
            if contracts and d["contract"] not in contracts:
                continue
            if d["ok"]:
                what, label = f"{d['type']} {describe(d['contract'])}", "deploy"
            else:
                what, label = d["type"], "deploy reverted"
            rows.append((d["block"], d["order"],
                         f"[{d['block']}] tx {h(d['tx'])} contract <{a(d['contract']) or 'none'}> {what} "
                         f"by <{a(d['signer'])}>  {label}"))

    for log, tx in zip(transfers, txs):
        contract = log["address"].lower()
        token_id = int(log["topics"][3], 16)
        label = labels.get((log["transactionHash"], contract, token_id))
        if label is None:
            label = ("mint" if log["topics"][1] == ZERO_TOPIC
                     else "burn" if log["topics"][2] == ZERO_TOPIC else "transfer")
        ctype = types.get(contract, "?")
        block = int(log["blockNumber"], 16)
        # A resale is always signed by the buyer, so "by" would repeat "to".
        by = "" if label.startswith("resale") else f" by <{a(tx['from'].lower())}>"
        rows.append((block, (int(log["transactionIndex"], 16), int(log["logIndex"], 16)),
                     f"[{block}] tx {h(log['transactionHash'])} contract <{a(contract)}> {ctype} {describe(contract)} "
                     f"tID {token_id} from <{a(word_to_address(log['topics'][1]))}> "
                     f"to <{a(word_to_address(log['topics'][2]))}>{by}  {label}"))

    for log, tx in zip(listings, listing_txs):
        contract = log["address"].lower()
        token_id = int(log["topics"][1], 16)
        signer = tx["from"].lower()
        if log["topics"][0] == LISTED_TOPIC:
            holder, label = word_to_address(log["topics"][2]), f"list price {int(log['data'], 16)}"
        else:
            holder, label = signer, "cancel listing"
        if args.sender and holder.lower() != args.sender.lower():
            continue
        ctype = types.get(contract, "?")
        block = int(log["blockNumber"], 16)
        # The holder's own wallet lists and cancels, so the seller is the signer.
        rows.append((block, (int(log["transactionIndex"], 16), int(log["logIndex"], 16)),
                     f"[{block}] tx {h(log['transactionHash'])} contract <{a(contract)}> {ctype} {describe(contract)} "
                     f"tID {token_id} seller <{a(holder)}>  {label}"))

    for _, _, line in sorted(rows):
        print(line)
    print(f"{len(rows)} found", file=sys.stderr)
    return 0


if __name__ == "__main__":
    sys.exit(main())
