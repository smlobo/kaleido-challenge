#!/usr/bin/env python3
"""Show an operation (GET /operations/{id}): how a transaction you started is doing.

Every write that touches the chain (creating an event, buying a ticket,
checking one in) answers with an operation ID. Only the caller who started
it can read it, so this needs that caller's token: --token, or else an
environment variable, ORG_TOKEN by default or the one --token-env names
(e.g. --token-env ALICE for a user's purchase). The variable keeps it out of
your shell history and the process list.

Statuses: signed (saved, not yet accepted by the node) -> submitted (waiting
to be mined) -> mined (succeeded), reverted (mined, but the contract said no;
the reason is shown) or failed (the node would never mine it).

Usage: scripts/get_operation.py OPERATION_ID [--wait] [--token TOKEN | --token-env VAR] [--url URL]
  --wait           poll until the operation is final (up to 60 s)
  --token TOKEN    the token of whoever started the operation
  --token-env VAR  environment variable holding the token, if no --token (default ORG_TOKEN)
  --url URL        the service, default $TICKETING_URL or http://localhost:8080
Exit status: 0 mined or still in progress, 1 reverted, failed or an error.
Standard library only.
"""

import argparse
import sys
import time

from ticketingapi import DEFAULT_URL, APIError, get_json, resolve_token

FINAL = {"mined", "reverted", "failed"}
WAIT_SECONDS = 60


def main() -> int:
    p = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    p.add_argument("operation_id")
    p.add_argument("--wait", action="store_true", help=f"poll until final (up to {WAIT_SECONDS} s)")
    p.add_argument("--token", metavar="TOKEN", help="the token of whoever started the operation")
    p.add_argument("--token-env", default="ORG_TOKEN", metavar="VAR",
                   help="environment variable holding the token, if no --token (default ORG_TOKEN)")
    p.add_argument("--url", default=DEFAULT_URL)
    args = p.parse_args()

    token = resolve_token(args.token, args.token_env)
    if not token:
        print(f"no token: pass --token or export {args.token_env} "
              "(the token of whoever started the operation)", file=sys.stderr)
        return 2

    path = "/operations/" + args.operation_id
    try:
        op = get_json(args.url, path, token=token)
        deadline = time.monotonic() + WAIT_SECONDS
        while args.wait and op["status"] not in FINAL and time.monotonic() < deadline:
            time.sleep(1)
            op = get_json(args.url, path, token=token)
    except APIError as e:
        print(e, file=sys.stderr)
        if "404" in str(e):
            print("(an operation is only visible to the caller who started it; check the token)", file=sys.stderr)
        return 1

    print(f"operation   {op['id']}")
    print(f"status      {op['status']}")
    print(f"from        {op['from_address']}  (nonce {op['nonce']})")
    print(f"tx hash     {op['tx_hash']}")
    if op.get("block_number") is not None:
        print(f"block       {op['block_number']}")
    if op.get("contract_address"):
        print(f"contract    {op['contract_address']}")
    if op.get("revert_reason"):
        print(f"reverted    {op['revert_reason']}")
    if op.get("last_error"):
        print(f"last error  {op['last_error']}  (after {op['attempts']} attempts)")
    sys.stdout.flush()

    if args.wait and op["status"] not in FINAL:
        print(f"Still {op['status']} after {WAIT_SECONDS} s; is Besu producing blocks?", file=sys.stderr)
    return 1 if op["status"] in ("reverted", "failed") else 0


if __name__ == "__main__":
    sys.exit(main())
