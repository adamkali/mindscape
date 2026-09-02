# Authentication & Session Design

*Plain-English explanation of how login sessions work in Mindscape, why the
design changed, and how passkeys will fit in later.*

## The problem, in everyday terms

Previously, logging in handed you a single "badge" (a JWT) that was valid for
3 days. Two problems:

- When the 3 days were up, the badge just stopped working — there was no way
  to renew it, so you were kicked back to the login screen.
- The server kept exactly **one** badge record per person. Log in on your
  laptop, then your phone — the phone's badge *replaced* the laptop's record,
  so the laptop got rejected on its next request.

## The fix: two tokens instead of one

Authentication is split into two pieces, which is the industry-standard
pattern:

- **Access token** — a short-lived badge (15 minutes, `server.access_token_ttl`).
  Sent with every API request in the `Authorization` header, exactly like
  before. Because it's so short-lived, a stolen one is nearly useless.
- **Refresh token** — a long-lived "renewal ticket" (30 days,
  `server.refresh_token_ttl`). Stored in an **httpOnly cookie**, which means
  JavaScript on the page *cannot read it* — even a malicious script injected
  into the site can't steal it. The browser only sends it to one URL:
  `/api/users/refresh`.

## How a normal day works

1. You log in. The server creates a **session record** for *this specific
   browser* (a row in the `sessions` table) and hands back both tokens.
2. You use the app. Every request carries the 15-minute access token.
3. After 15 minutes the access token expires. The next API call gets a 401 —
   but instead of logging you out, the frontend quietly calls
   `POST /api/users/refresh`. The browser attaches the refresh cookie, the
   server checks it against your session record, issues a fresh access token,
   and **rotates** the refresh token (old ticket destroyed, new one issued).
   Your original request is retried. You never notice.
4. This continues for as long as you keep using the app. You only see the
   login screen again if you're away long enough for the refresh token itself
   to expire (30 days), or you explicitly log out.

## Why two browsers can now coexist

Each login creates its **own** session row — laptop and phone each have their
own renewal ticket. Logging in somewhere new *adds* a row; it never
overwrites another. Logging out (`DELETE /api/users/refresh`) deletes only
that one row. This also makes a future "log out everywhere" button trivial:
delete all rows for the user (`DeleteSessionsByUserId` already exists).

## Security details, briefly

- The server never stores the refresh token itself — only a SHA-256
  fingerprint (hash) in `sessions.refresh_token_hash`. A database leak
  doesn't leak usable tickets.
- Refresh tokens are **rotated on every use**. If two requests race to
  refresh at once, exactly one wins (a single atomic `UPDATE ... WHERE
  refresh_token_hash = old`); the loser gets a 401 and recovers via the
  rotated cookie. The frontend also single-flights refreshes
  (`web/src/utils/refreshAuth.ts`), so this race is rare in practice.
- The cookie is `HttpOnly` (JS can't read it), `Secure` (HTTPS only — set
  `server.cookie.secure: true` in production), `SameSite=Strict` (other
  websites can't trigger it), and scoped to `Path=/api/users/refresh` (sent
  nowhere else). Logout is `DELETE` on the same path so the cookie reaches it.
- Access-token checks are **stateless** — no database lookup per request
  (`AuthService.CheckToken` verifies signature + expiry only). Revocation
  happens at the refresh layer: kill the session row and the user is out
  within one access-token TTL (15 minutes) at most.

## Key code

| Piece | Location |
|---|---|
| Session table | `db/migrations/20260604000000_create_sessions.sql` |
| Session queries | `db/queries/session.sql` |
| Issue / refresh / revoke | `services/AuthService.go` (`IssueSession`, `RefreshSession`, `RevokeSession`) |
| Cookie helpers | `models/handlers/user_handlers/session_cookie.go` |
| Refresh / logout endpoints | `POST` / `DELETE /api/users/refresh` (`controllers/user_controller.go`) |
| Silent refresh (frontend) | `web/src/utils/refreshAuth.ts` + `web/src/utils/authInterceptor.ts` |
| TTL config | `server.access_token_ttl`, `server.refresh_token_ttl`, `server.cookie.secure` |

The legacy `tokens` table is no longer written to; it is kept for now so the
migration is reversible.

## What about passkeys?

This design is the foundation. Login, refresh, and (later) passkey login all
funnel into one shared helper, `AuthService.IssueSession(user, userAgent)`,
which mints the token pair and session row. Adding passkeys later means:
verify the WebAuthn assertion (via `github.com/go-webauthn/webauthn`), then
call `IssueSession` — no session rework needed.
