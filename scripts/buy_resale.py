#!/usr/bin/env python3
"""Buy a ticket listed for resale (POST /tickets/{id}/purchase) and wait until it is yours.

The purchase is a transaction from your own custodial wallet. --max-price is
the most you agree to pay; without it, the price listed now. If the seller
raises the price before the purchase is mined, the contract refuses it. No
money moves yet: the contract records the price.

Needs the buyer's API token: --token, or else an environment variable,
USER_TOKEN by default or the one --token-env names (e.g. --token-env BOB).
To find tickets for sale: curl localhost:8080/events/<event-id>/listings

Usage: scripts/buy_resale.py TICKET_ID [options]
  --max-price N       the most you agree to pay, in base units (default: the listed price)
  --token TOKEN       the buyer's token
  --token-env VAR     environment variable holding the token, if no --token (default USER_TOKEN)
  --idempotency-key K reuse to retry the same purchase safely (default: a new UUID)
  --no-wait           print as soon as the purchase is saved
  --url URL           the service, default $TICKETING_URL or http://localhost:8080
Exit status: 0 bought (or saved, with --no-wait), 1 refused, reverted or an error, 2 no token.
Standard library only.
"""

import argparse
import sys
import uuid

from ticketingapi import DEFAULT_URL, APIError, get_json, post_json, print_outcome, resolve_token, wait_for_operation

WAIT_SECONDS = 30


def main() -> int:
    p = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    p.add_argument("ticket_id")
    p.add_argument("--max-price", metavar="N")
    p.add_argument("--token", metavar="TOKEN", help="the buyer's token (default: the --token-env variable)")
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
    body = {"max_price": args.max_price} if args.max_price else {}

    op = None
    try:
        action = post_json(args.url, f"/tickets/{args.ticket_id}/purchase", body, token=token,
                           headers={"Idempotency-Key": key})
        if not args.no_wait:
            op = wait_for_operation(args.url, action["operation_id"], token, WAIT_SECONDS)
            ticket = get_json(args.url, "/tickets/" + args.ticket_id)
        else:
            ticket = action["ticket"]
    except APIError as e:
        print(e, file=sys.stderr)
        if e.retryable:
            print(f"(to retry this same purchase: --idempotency-key {key})", file=sys.stderr)
        return 1

    print(f"ticket      {ticket['id']}")
    print(f"number      {ticket['token_id']}")
    print(f"owner       {ticket.get('owner_address')}")
    print(f"operation   {action['operation_id']}")
    sys.stdout.flush()
    return print_outcome(op, WAIT_SECONDS)


if __name__ == "__main__":
    sys.exit(main())
