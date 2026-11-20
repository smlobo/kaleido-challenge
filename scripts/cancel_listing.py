#!/usr/bin/env python3
"""Take your ticket off the resale market (DELETE /tickets/{id}/listing) and wait until it is done.

The cancellation is a transaction from your own custodial wallet: the
contract checks you hold the ticket and that it is listed. By default this
waits until the transaction is mined.

Needs the holder's API token: --token, or else an environment variable,
USER_TOKEN by default or the one --token-env names (e.g. --token-env ALICE).

Usage: scripts/cancel_listing.py TICKET_ID [options]
  --token TOKEN       the holder's token
  --token-env VAR     environment variable holding the token, if no --token (default USER_TOKEN)
  --idempotency-key K reuse to retry the same cancellation safely (default: a new UUID)
  --no-wait           print as soon as the cancellation is saved
  --url URL           the service, default $TICKETING_URL or http://localhost:8080
Exit status: 0 cancelled (or saved, with --no-wait), 1 refused, reverted or an error, 2 no token.
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
        action = request("DELETE", args.url, f"/tickets/{args.ticket_id}/listing", token=token,
                         headers={"Idempotency-Key": key})
        if not args.no_wait:
            op = wait_for_operation(args.url, action["operation_id"], token, WAIT_SECONDS)
            ticket = get_json(args.url, "/tickets/" + args.ticket_id)
        else:
            ticket = action["ticket"]
    except APIError as e:
        print(e, file=sys.stderr)
        if e.retryable:
            print(f"(to retry this same cancellation: --idempotency-key {key})", file=sys.stderr)
        return 1

    print(f"ticket      {ticket['id']}")
    print(f"number      {ticket['token_id']}")
    print(f"listed at   {ticket.get('listing_price') or '(not listed)'}")
    print(f"operation   {action['operation_id']}")
    sys.stdout.flush()
    return print_outcome(op, WAIT_SECONDS)


if __name__ == "__main__":
    sys.exit(main())
