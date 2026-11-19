#!/usr/bin/env python3
"""Create a user on the ticketing service (POST /users) and print the result.

The service creates the user's custodial wallet and an API token. The token is
shown only in this response, so keep it: it is what the user sends as
"Authorization: Bearer <token>" to buy tickets.

Usage: scripts/create_user.py NAME [--url URL] [--token-only]
  --url         the service, default $TICKETING_URL or http://localhost:8080
  --token-only  print just the token, for use in a shell:
                ALICE=$(scripts/create_user.py alice --token-only)
Standard library only.
"""

import argparse
import sys

from ticketingapi import DEFAULT_URL, APIError, post_json


def main() -> int:
    p = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    p.add_argument("name", help="the user's login name (subject); must be unused")
    p.add_argument("--url", default=DEFAULT_URL, help=f"the service (default {DEFAULT_URL})")
    p.add_argument("--token-only", action="store_true", help="print only the API token")
    args = p.parse_args()

    try:
        user = post_json(args.url, "/users", {"subject": args.name})
    except APIError as e:
        print(e, file=sys.stderr)
        return 1

    if args.token_only:
        print(user["token"])
    else:
        print(f"user    {user['id']}")
        print(f"name    {user['subject']}")
        print(f"wallet  {user['wallet_address']}")
        print(f"token   {user['token']}", flush=True)
        print("Keep the token: the service stores only its hash and cannot show it again.", file=sys.stderr)
    return 0


if __name__ == "__main__":
    sys.exit(main())
