#!/usr/bin/env python3
"""Create an organizer on the ticketing service (POST /organizers) and print the result.

Needs the platform admin token (the same value the service was started
with): --token, or else the ADMIN_TOKEN environment variable. The variable
keeps it out of your shell history and the process list.

The service creates the organizer's box office, gate and payout keys and an
API token. The token is shown only in this response, so keep it: it is what
the organizer sends as "Authorization: Bearer <token>" to create events and
check tickets in.

Usage: scripts/create_organizer.py NAME [--token TOKEN] [--url URL] [--token-only]
  --token       the admin token (default $ADMIN_TOKEN)
  --url         the service, default $TICKETING_URL or http://localhost:8080
  --token-only  print just the token, for use in a shell:
                ORG=$(scripts/create_organizer.py "Rock Fest Ltd" --token-only)
Standard library only.
"""

import argparse
import sys

from ticketingapi import DEFAULT_URL, APIError, add_token_arg, post_json, resolve_token


def main() -> int:
    p = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    p.add_argument("name", help="the organizer's name, e.g. \"Rock Fest Ltd\"")
    p.add_argument("--url", default=DEFAULT_URL, help=f"the service (default {DEFAULT_URL})")
    p.add_argument("--token-only", action="store_true", help="print only the API token")
    add_token_arg(p, "ADMIN_TOKEN", "the platform admin")
    args = p.parse_args()

    admin = resolve_token(args.token, "ADMIN_TOKEN")
    if not admin:
        print("no admin token: pass --token or export ADMIN_TOKEN (the value the service was started with)",
              file=sys.stderr)
        return 2

    try:
        org = post_json(args.url, "/organizers", {"name": args.name}, token=admin)
    except APIError as e:
        print(e, file=sys.stderr)
        return 1

    if args.token_only:
        print(org["token"])
    else:
        print(f"organizer   {org['id']}")
        print(f"name        {org['name']}")
        print(f"box office  {org['box_office_address']}  (mints tickets)")
        print(f"gate        {org['gate_address']}  (checks tickets in)")
        print(f"payout      {org['payout_address']}  (will receive payments)")
        print(f"token       {org['token']}", flush=True)
        print("Keep the token: the service stores only its hash and cannot show it again.", file=sys.stderr)
    return 0


if __name__ == "__main__":
    sys.exit(main())
