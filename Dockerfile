# Build: docker build -f Dockerfile -t patriot-pest-go .
# Production image for patriotpest.pro. The SQLite database is NOT baked in —
# it is bind-mounted at runtime onto /app/database/patriot.db.
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux GOTOOLCHAIN=local go build -ldflags "-X main.version=${VERSION}" -o /out/patriot-server ./cmd/server \
 && mkdir -p /out/storage/logs /out/database /out/geo && chmod 750 /out/storage/logs /out/database
# Coarse IP geolocation database for analytics. Best effort: if the
# download fails (offline build) the build still succeeds and the server
# simply reports country/city as unknown. Raw IPs are never stored.
# Uses busybox wget (always present in alpine); verifies a plausible size.
RUN (wget -q -O /out/geo/GeoLite2-City.mmdb \
	"https://github.com/P3TERX/GeoLite.mmdb/raw/download/GeoLite2-City.mmdb" \
	&& [ $(stat -c%s /out/geo/GeoLite2-City.mmdb) -gt 10000000 ] \
	|| { echo "geo db download failed, continuing without it"; rm -f /out/geo/GeoLite2-City.mmdb; })

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /out/patriot-server /app/patriot-server
COPY --from=build /out/geo /app/geo
COPY --from=build /src/migrations /app/migrations
COPY --from=build /src/configs /app/configs
COPY --from=build --chown=65532:65532 --chmod=750 /out/database /app/database
COPY --from=build /out/storage /app/storage
EXPOSE 3000
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 CMD ["/app/patriot-server","-health"]
ENTRYPOINT ["/app/patriot-server"]
