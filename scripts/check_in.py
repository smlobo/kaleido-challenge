#!/usr/bin/env python3
"""Check a ticket in at the gate (POST /tickets/{id}/check-in) and wait until it is done.

Needs the token of the organizer whose event the ticket is for: --token, or
else the ORG_TOKEN environment variable, which keeps it out of your shell
history and the process list. The organizer's gate key burns the ticket
on-chain, so it can never be used again.

The service answers at once; the burn is mined a block or two later. By
default this script waits for that and prints the result, with the
contract's reason if it refused.

Usage: scripts/check_in.py TICKET_ID [options]
  --token TOKEN       the organizer's token (default $ORG_TOKEN)
  --idempotency-key K reuse to retry a check-in safely (default: a new UUID)
  --no-wait           print as soon as the check-in is saved
  --url URL           the service, default $TICKETING_URL or http://localhost:8080
Exit status: 0 checked in (or saved, with --no-wait), 1 refused, reverted or
an error, 2 no token.
Standard library only.
"""

import argparse
import sys
import time
import uuid

from ticketingapi import DEFAULT_URL, APIError, add_token_arg, get_json, post_json, resolve_token

FINAL = {"mined", "reverted", "failed"}
WAIT_SECONDS = 30


def main() -> int:
    p = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    p.add_argument("ticket_id")
    add_token_arg(p, "ORG_TOKEN", "the event organizer")
    p.add_argument("--idempotency-key", default=None)
    p.add_argument("--no-wait", action="store_true")
    p.add_argument("--url", default=DEFAULT_URL)
    args = p.parse_args()

    token = resolve_token(args.token, "ORG_TOKEN")
    if not token:
        print("no organizer token: pass --token or export ORG_TOKEN (from create_organizer.py)", file=sys.stderr)
        return 2
    key = args.idempotency_key or str(uuid.uuid4())

    op = None
    try:
        ticket = post_json(args.url, f"/tickets/{args.ticket_id}/check-in", {}, token=token,
                           headers={"Idempotency-Key": key})
        op_id = ticket["checkin_operation_id"]
        if not args.no_wait:
            # Follow the operation rather than the ticket: if the contract
            # refuses, the operation says why, while the ticket just stays active.
            deadline = time.monotonic() + WAIT_SECONDS
            op = get_json(args.url, "/operations/" + op_id, token=token)
            while op["status"] not in FINAL and time.monotonic() < deadline:
                time.sleep(1)
                op = get_json(args.url, "/operations/" + op_id, token=token)
            ticket = get_json(args.url, "/tickets/" + ticket["id"])
    except APIError as e:
        print(e, file=sys.stderr)
        if e.retryable:
            print(f"(to retry this same check-in: --idempotency-key {key})", file=sys.stderr)
        return 1

    print(f"ticket      {ticket['id']}")
    print(f"number      {ticket['token_id']}")
    print(f"status      {ticket['status']}")
    print(f"operation   {op_id}")
    if op is not None:
        print(f"tx          {op['status']}" + (f" in block {op['block_number']}" if op.get("block_number") else ""))
        if op.get("revert_reason"):
            print(f"reverted    {op['revert_reason']}")
        if op.get("last_error"):
            print(f"last error  {op['last_error']}")
    sys.stdout.flush()

    if op is None:
        return 0
    if op["status"] in ("reverted", "failed"):
        print("The check-in did not happen; the ticket is still valid and can be checked in again.", file=sys.stderr)
        return 1
    if op["status"] not in FINAL:
        print(f"Still {op['status']} after {WAIT_SECONDS} s; is Besu producing blocks?", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
