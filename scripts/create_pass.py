#!/usr/bin/env python3
"""Make a check-in pass for your ticket (POST /tickets/{id}/pass): what you show at the door.

Check-in needs both parties on-chain: the organizer's gate sends the
transaction, and it carries this pass, which the ticket's holder signs. The
service signs it with the holder's custodial wallet key; nothing is sent to
the chain. A pass is valid for five minutes and works once.

Needs the holder's API token (from create_user.py): --token, or else an
environment variable, USER_TOKEN by default or the one --token-env names
(e.g. --token-env ALICE). The variable keeps it out of your shell history and
the process list.

Prints the pass as one line of JSON (what a QR code at the door would hold),
for the organizer to pass to check_in.py:
  PASS=$(scripts/create_pass.py $TICKET)
  scripts/check_in.py $TICKET --pass "$PASS"

Usage: scripts/create_pass.py TICKET_ID [options]
  --token TOKEN    the holder's token
  --token-env VAR  environment variable holding the token, if no --token (default USER_TOKEN)
  --url URL        the service, default $TICKETING_URL or http://localhost:8080
Exit status: 0 made, 1 refused or an error, 2 no token.
Standard library only.
"""

import argparse
import json
import sys

from ticketingapi import DEFAULT_URL, APIError, request, resolve_token


def main() -> int:
    p = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    p.add_argument("ticket_id")
    p.add_argument("--token", metavar="TOKEN", help="the holder's token (default: the --token-env variable)")
    p.add_argument("--token-env", default="USER_TOKEN", metavar="VAR")
    p.add_argument("--url", default=DEFAULT_URL)
    args = p.parse_args()

    token = resolve_token(args.token, args.token_env)
    if not token:
        print(f"no user token: pass --token or export {args.token_env} (from create_user.py)", file=sys.stderr)
        return 2

    try:
        ticket_pass = request("POST", args.url, f"/tickets/{args.ticket_id}/pass", token=token)
    except APIError as e:
        print(e, file=sys.stderr)
        return 1

    print(json.dumps(ticket_pass, separators=(",", ":")))
    print(f"pass for ticket #{ticket_pass['token_id']}, valid until {ticket_pass['expires_at']}", file=sys.stderr)
    return 0


if __name__ == "__main__":
    sys.exit(main())
