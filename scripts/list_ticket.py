#!/usr/bin/env python3
"""List your ticket for resale (PUT /tickets/{id}/listing), or change its price.

The listing is a transaction from your own custodial wallet: the contract
checks you hold the ticket, that the price is no more than the event's
resale cap (face value x max_resale_bps / 10000), and that the event has not
started. Listing again changes the price. By default this waits until the
transaction is mined. To take the ticket off the market: scripts/cancel_listing.py

Needs the holder's API token: --token, or else an environment variable,
USER_TOKEN by default or the one --token-env names (e.g. --token-env ALICE).

Usage: scripts/list_ticket.py TICKET_ID --price N [options]
  --price N           asking price in base units (e.g. 110000000 = 110.00 in a 6-decimal currency)
  --token TOKEN       the holder's token
  --token-env VAR     environment variable holding the token, if no --token (default USER_TOKEN)
  --idempotency-key K reuse to retry the same request safely (default: a new UUID)
  --no-wait           print as soon as the request is saved
  --url URL           the service, default $TICKETING_URL or http://localhost:8080
Exit status: 0 done (or saved, with --no-wait), 1 refused, reverted or an error, 2 no token.
Standard library only.
"""

import argparse
import sys
import uuid

from ticketingapi import DEFAULT_URL, APIError, get_json, print_outcome, request, resolve_token, wait_for_operation

WAIT_SECONDS = 30


def main() -> int:
    p = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    p.add_argument("ticket_id")
    p.add_argument("--price", metavar="N", required=True)
    p.add_argument("--token", metavar="TOKEN", help="the holder's token (default: the --token-env variable)")
    p.add_argument("--token-env", default="USER_TOKEN", metavar="VAR")
    p.add_argument("--idempotency-key", default=None)
    p.add_argument("--no-wait", action="store_true")
    p.add_argument("--url", default=DEFAULT_URL)
    args = p.parse_args()

    token = resolve_token(args.token, args.token_env)
    if not token:
        print(f"no user token: pass --token or export {args.token_env} (from create_user.py)", file=sys.stderr)
        return 2
    key = args.idempotency_key or str(uuid.uuid4())

    op = None
    try:
        action = request("PUT", args.url, f"/tickets/{args.ticket_id}/listing", {"price": args.price}, token=token,
                         headers={"Idempotency-Key": key})
        if not args.no_wait:
            op = wait_for_operation(args.url, action["operation_id"], token, WAIT_SECONDS)
            ticket = get_json(args.url, "/tickets/" + args.ticket_id)
        else:
            ticket = action["ticket"]
    except APIError as e:
        print(e, file=sys.stderr)
        if e.retryable:
            print(f"(to retry this same request: --idempotency-key {key})", file=sys.stderr)
        return 1

    print(f"ticket      {ticket['id']}")
    print(f"number      {ticket['token_id']}")
    print(f"listed at   {ticket.get('listing_price') or '(not listed)'}")
    print(f"operation   {action['operation_id']}")
    sys.stdout.flush()
    return print_outcome(op, WAIT_SECONDS)


if __name__ == "__main__":
    sys.exit(main())
