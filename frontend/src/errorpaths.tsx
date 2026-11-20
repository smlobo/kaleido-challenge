import { Card } from "./ui";

// What each refused action returns. "Service" means the service's own checks
// refused it; "contract" means the contract's dry run (eth_estimateGas)
// reverted, and the service decoded the revert. Either way, nothing was signed.
const PATHS: { try: string; where: string; status: number; code: string; who: "service" | "contract" | "API layer" }[] = [
  { try: "Buy for an event id that doesn't exist (paste any uuid)", where: "Tickets: buy", status: 404, code: "not_found", who: "service" },
  { try: "Buy for an event id that isn't a uuid", where: "Tickets: buy", status: 400, code: "invalid_request", who: "API layer" },
  { try: "Buy for an event still deploying", where: "Tickets: buy", status: 409, code: "event_not_ready", who: "service" },
  { try: "Buy for a sold-out event (create one with supply 1, buy twice)", where: "Tickets: buy", status: 422, code: "sold_out", who: "service" },
  { try: "Buy with an organizer's token", where: "Tickets: buy", status: 403, code: "forbidden", who: "service" },
  { try: "Buy with a pasted token that isn't one", where: "Tickets: buy", status: 401, code: "unauthorized", who: "API layer" },
  { try: "Buy with no token", where: "Tickets: buy", status: 401, code: "unauthorized", who: "service" },
  { try: "Buy your own listed ticket (act as the owner, Buy)", where: "Tickets", status: 422, code: "buyer_is_seller", who: "service" },
  { try: "List, Cancel or Buy a checked-in ticket", where: "Tickets", status: 409, code: "already_checked_in", who: "service" },
  { try: "Show pass for a checked-in ticket", where: "Tickets", status: 409, code: "already_checked_in", who: "service" },
  { try: "List, Cancel or Show pass for someone else's ticket", where: "Tickets", status: 403, code: "not_your_ticket", who: "service" },
  { try: "Cancel or Buy a ticket that isn't listed", where: "Tickets", status: 409, code: "not_listed", who: "service" },
  { try: "List above the resale cap", where: "Tickets", status: 422, code: "price_above_cap", who: "service" },
  { try: "List at price 0", where: "Tickets", status: 400, code: "invalid_price", who: "service" },
  { try: "Buy with a max price below the listing (as if the seller raised it)", where: "Tickets", status: 422, code: "price_changed", who: "service" },
  { try: "List a ticket whose mint is still in flight", where: "Tickets", status: 409, code: "not_minted_yet", who: "service" },
  { try: "Any action while the ticket's last transaction is still in flight", where: "Tickets", status: 409, code: "action_in_progress", who: "service" },
  { try: "Check in the same pass twice (replay)", where: "Gate", status: 409, code: "already_checked_in", who: "service" },
  { try: "Expire the pass", where: "Gate", status: 422, code: "pass_expired", who: "service" },
  { try: "Extend the deadline, or forge the signature", where: "Gate", status: 422, code: "pass_not_from_holder", who: "contract" },
  { try: "Seller's pass after a resale: Show pass as the seller, hand it to the gate, resell, then check in", where: "Gate", status: 422, code: "pass_not_from_holder", who: "contract" },
  { try: "Check in at another organizer's gate", where: "Gate", status: 403, code: "not_your_event", who: "service" },
];

export function ErrorPaths() {
  return (
    <Card
      step="!"
      title="Error paths to try"
      intro={
        <>
          Every action is now offered, whatever the ticket's state, and each ticket can be acted on as any user. The
          service refuses before anything is signed, with an RFC 9457 problem. Where the contract would refuse, the
          service finds out from a dry run (<code>eth_estimateGas</code>) and decodes the revert.
        </>
      }
    >
      <div className="table-wrap">
        <table>
          <thead>
            <tr>
              <th>Try</th>
              <th>Where</th>
              <th>Expect</th>
              <th>Refused by</th>
            </tr>
          </thead>
          <tbody>
            {PATHS.map((p) => (
              <tr key={p.try}>
                <td className="wrap">{p.try}</td>
                <td>{p.where}</td>
                <td>
                  <code>
                    {p.status} {p.code}
                  </code>
                </td>
                <td>{p.who}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </Card>
  );
}
