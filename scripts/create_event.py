#!/usr/bin/env python3
"""Create an event on the ticketing service (POST /events) and wait until it is on sale.

Needs an organizer's API token (the one create_organizer.py printed): --token,
or else the ORG_TOKEN environment variable, which keeps it out of your shell
history and the process list. The event belongs to
that organizer: its box office key will mint the tickets and its gate key
will check them in.

The service answers at once with the event in "deploying"; its TicketNFT
contract is mined a block or two later. By default this script waits for
that and prints the contract address.

Usage: scripts/create_event.py NAME [options]
  --token TOKEN       the organizer's token (default $ORG_TOKEN)
  --symbol SYM        token symbol, 1-11 capital letters or digits
                      (default: the name's letters and digits, capitalized and
                      cut to 11, e.g. "Rock Fest 2027" -> "ROCKFEST202")
  --starts-in DUR     when the event starts, from now: 90m, 6h, 2d (default 1d)
  --supply N          number of tickets (default 100)
  --face-value N      price in base units, a whole number (default 100000000,
                      i.e. 100.00 in a 6-decimal currency)
  --max-resale-bps N  resale cap in basis points of face value (default 11000 = 110%)
  --idempotency-key K reuse to retry a request safely (default: a new UUID)
  --no-wait           print the event as soon as it is saved
  --url URL           the service, default $TICKETING_URL or http://localhost:8080
  --id-only           print just the event ID, for use in a shell:
                      EVENT=$(scripts/create_event.py "Rock Fest 2027" --id-only)
Standard library only.
"""

import argparse
import re
import sys
import time
import uuid
from datetime import datetime, timedelta, timezone

from ticketingapi import DEFAULT_URL, APIError, add_token_arg, get_json, post_json, resolve_token

WAIT_SECONDS = 30


def duration(s: str) -> timedelta:
    m = re.fullmatch(r"(\d+)([mhd])", s.strip())
    if not m:
        raise argparse.ArgumentTypeError(f"{s!r}: use a number and m, h or d, e.g. 90m, 6h, 2d")
    n, unit = int(m.group(1)), m.group(2)
    return {"m": timedelta(minutes=n), "h": timedelta(hours=n), "d": timedelta(days=n)}[unit]


def default_symbol(name: str) -> str:
    return re.sub(r"[^A-Z0-9]", "", name.upper())[:11] or "EVENT"


def main() -> int:
    p = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    p.add_argument("name", help='the event\'s name, e.g. "Rock Fest 2027"')
    p.add_argument("--symbol")
    p.add_argument("--starts-in", type=duration, default=timedelta(days=1), metavar="DUR")
    p.add_argument("--supply", type=int, default=100)
    p.add_argument("--face-value", default="100000000")
    p.add_argument("--max-resale-bps", type=int, default=11_000)
    p.add_argument("--idempotency-key", default=None)
    p.add_argument("--no-wait", action="store_true")
    p.add_argument("--url", default=DEFAULT_URL)
    p.add_argument("--id-only", action="store_true")
    add_token_arg(p, "ORG_TOKEN", "the organizer")
    args = p.parse_args()

    token = resolve_token(args.token, "ORG_TOKEN")
    if not token:
        print("no organizer token: pass --token or export ORG_TOKEN (from create_organizer.py)", file=sys.stderr)
        return 2

    starts_at = (datetime.now(timezone.utc) + args.starts_in).replace(microsecond=0)
    body = {
        "name": args.name,
        "symbol": args.symbol or default_symbol(args.name),
        "starts_at": starts_at.isoformat().replace("+00:00", "Z"),
        "max_supply": args.supply,
        "face_value": args.face_value,
        "max_resale_bps": args.max_resale_bps,
    }
    key = args.idempotency_key or str(uuid.uuid4())

    try:
        event = post_json(args.url, "/events", body, token=token, headers={"Idempotency-Key": key})
        if not args.no_wait:
            deadline = time.monotonic() + WAIT_SECONDS
            while event["status"] == "deploying" and time.monotonic() < deadline:
                time.sleep(1)
                event = get_json(args.url, "/events/" + event["id"])
    except APIError as e:
        print(e, file=sys.stderr)
        if e.retryable:
            print(f"(to retry this same request: --idempotency-key {key})", file=sys.stderr)
        return 1

    if args.id_only:
        print(event["id"])
    else:
        print(f"event       {event['id']}")
        print(f"name        {event['name']} ({event['symbol']})")
        print(f"starts      {event['starts_at']}")
        print(f"tickets     {event['max_supply']} at face value {event['face_value']}, "
              f"resale capped at {event['max_resale_bps'] / 100:g}%")
        print(f"status      {event['status']}")
        print(f"contract    {event.get('contract_address') or '(not deployed yet)'}")
        print(f"operation   {event['deploy_operation_id']}", flush=True)

    if event["status"] == "failed":
        print("The deploy failed; GET /operations/<operation> with ORG_TOKEN shows why.", file=sys.stderr)
        return 1
    if event["status"] == "deploying" and not args.no_wait:
        print(f"Still deploying after {WAIT_SECONDS} s; is Besu producing blocks?", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
