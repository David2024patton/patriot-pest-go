# patriot-pest-go

The public landing site for **patriotpest.pro**: marketing pages, the content
catalog (pest library, service areas, blog), lead capture, and the PWA shell.

## What is NOT here

Customer, staff and admin dashboards — and every FieldRoutes / Twilio
integration — live in the **AlphaFlux platform** (`alphaflux.net`, repo
`David2024patton/alphaflux-console`), which serves many companies, not just
Patriot Pest. This site does not authenticate anyone: every sign-in link, and
every URL that used to serve a dashboard, redirects to the platform

- `/login`, `/logout`, `/account`, `/dashboard`, `/customer-dashboard`,
  `/customer-portal`, `/staff-dashboard`, `/staff`, `/su` → `ALPHAFLUX_LOGIN_URL`
- `/admin`, `/admin/*` → same
- `/api/*` → `410 Gone` (the platform owns the API)

## Layout
- `cmd/server` — entrypoint: routes, redirects, graceful shutdown
- `internal/modules/marketing` — every public page, `/search`, sitemap, RSS, signup, contact, first-party beacon sink
- `internal/modules/health` — `/health`, `/ready`, `/metrics`
- `internal/view` — page templates + `internal/view/assets` (embedded, 7 MB)
- `internal/data` — SQLite content catalog loader
- `Dockerfile` — production image (**no database baked in**; the SQLite catalog is bind-mounted)

## Configuration
| Env | Default | Purpose |
| --- | --- | --- |
| `ADDR` | `:3000` | listen address |
| `DB_PATH` | `database/patriot.db` | SQLite content catalog |
| `APP_URL` | `https://patriotpest.pro` | canonical origin |
| `ALPHAFLUX_LOGIN_URL` | `https://alphaflux.net/login` | where sign-in links go |
| `MARKETING_ENABLED` | `true` | gates the public pages module |

## Run / build
```bash
go build -o /tmp/patriot-landing ./cmd/server
DB_PATH=database/patriot.db ADDR=:3120 APP_ENV=local /tmp/patriot-landing
```

## Production
Deployed by Dokploy (project **PPC** → app **Patriot**, service `ppc-patriot-qyfetu`)
from this repo, Dockerfile path `Dockerfile`. **Push to `main` auto-deploys.**
Runtime mounts:
- `/home/server/go-patriot/data/patriot.db` → `/app/database/patriot.db`
- `/home/server/go-patriot/data/settings.json` → `/app/storage/settings.json`

Domains: `patriotpest.pro`, `www.patriotpest.pro` (Let's Encrypt via the Dokploy Traefik).
