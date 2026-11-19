#!/usr/bin/env python3
"""Buy a ticket for an event (POST /events/{id}/tickets) and wait until it is minted.

Needs a user's API token (from create_user.py): --token, or else an
environment variable, USER_TOKEN by default or the one --token-env names
(e.g. --token-env ALICE). The variable keeps it out of your shell history and
the process list. The ticket is minted to that user's wallet by
the event organizer's box office key.

The service answers at once with the ticket in "pending_mint"; the mint is
mined a block or two later. By default this script waits for that and prints
the ticket's number and owner.

Usage: scripts/buy_ticket.py EVENT_ID [options]
  --token TOKEN       the user's token
  --token-env VAR     environment variable holding the user's token, if no --token
                      (default USER_TOKEN)
  --idempotency-key K reuse to retry a purchase safely (default: a new UUID)
  --no-wait           print the ticket as soon as it is saved
  --url URL           the service, default $TICKETING_URL or http://localhost:8080
  --id-only           print just the ticket ID, for use in a shell:
                      TICKET=$(scripts/buy_ticket.py $EVENT --token-env ALICE --id-only)
Standard library only.
"""

import argparse
import sys
import time
import uuid

from ticketingapi import DEFAULT_URL, APIError, get_json, post_json, resolve_token

WAIT_SECONDS = 30


def main() -> int:
    p = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    p.add_argument("event_id")
    p.add_argument("--token", metavar="TOKEN", help="the user's token (default: the --token-env variable)")
    p.add_argument("--token-env", default="USER_TOKEN", metavar="VAR")
    p.add_argument("--idempotency-key", default=None)
    p.add_argument("--no-wait", action="store_true")
    p.add_argument("--url", default=DEFAULT_URL)
    p.add_argument("--id-only", action="store_true")
    args = p.parse_args()

    token = resolve_token(args.token, args.token_env)
    if not token:
        print(f"no user token: pass --token or export {args.token_env} (from create_user.py)", file=sys.stderr)
        return 2
    key = args.idempotency_key or str(uuid.uuid4())

    try:
        ticket = post_json(args.url, f"/events/{args.event_id}/tickets", {}, token=token,
                           headers={"Idempotency-Key": key})
        if not args.no_wait:
            deadline = time.monotonic() + WAIT_SECONDS
            while ticket["status"] == "pending_mint" and time.monotonic() < deadline:
                time.sleep(1)
                ticket = get_json(args.url, "/tickets/" + ticket["id"])
    except APIError as e:
        print(e, file=sys.stderr)
        if e.retryable:
            print(f"(to retry this same purchase: --idempotency-key {key})", file=sys.stderr)
        return 1

    if args.id_only:
        print(ticket["id"])
    else:
        print(f"ticket      {ticket['id']}")
        print(f"event       {ticket['event_id']}")
        print(f"number      {ticket['token_id']}  (the token ID on-chain)")
        print(f"status      {ticket['status']}")
        print(f"owner       {ticket.get('owner_address') or '(set once the mint is mined)'}")
        print(f"record hash {ticket['record_hash']}")
        print(f"metadata    {ticket['metadata_url']}")
        print(f"operation   {ticket['mint_operation_id']}", flush=True)

    if ticket["status"] == "failed":
        print(f"The mint failed; scripts/get_operation.py {ticket['mint_operation_id']} "
              f"--token-env {args.token_env} shows why.", file=sys.stderr)
        return 1
    if ticket["status"] == "pending_mint" and not args.no_wait:
        print(f"Still pending after {WAIT_SECONDS} s; is Besu producing blocks?", file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    sys.exit(main())
