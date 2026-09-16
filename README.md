# patriot-pest-go

Go rewrite of the Patriot Pest Control site and dashboards (patriotpest.pro).

## Layout
- `cmd/server` — entrypoint
- `internal/` — modules (auth, portal, staffdash, admin, customers, fieldroutes, twilio, api, view/…)
- `pkg/` — shared packages
- `configs/`, `migrations/` — config defaults and SQL migrations
- `Dockerfile.go` — production image (**no database baked in**; the SQLite DB is bind-mounted at runtime)

## Build / run locally
```bash
go build -o /tmp/patriot-server ./cmd/server
OTP_DEBUG=1 ADDR=:3120 APP_ENV=local /tmp/patriot-server
```

## Production
Deployed by Dokploy (project **PPC** → app **Patriot**, service `ppc-patriot-qyfetu`) from this repo,
Dockerfile path `Dockerfile.go`. Runtime mounts:
- `/home/server/go-patriot/data/patriot.db` → `/app/database/patriot.db`
- `/home/server/go-patriot/data/settings.json` → `/app/storage/settings.json`

Domains: `patriotpest.pro`, `www.patriotpest.pro` (Let's Encrypt via the Dokploy Traefik).
