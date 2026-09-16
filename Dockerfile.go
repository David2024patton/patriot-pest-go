# Build: docker build -f Dockerfile.go -t patriot-pest-go .
# Production image for patriotpest.pro. The SQLite database is NOT baked in —
# it is bind-mounted at runtime onto /app/database/patriot.db.
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOTOOLCHAIN=local go build -o /out/patriot-server ./cmd/server \
 && mkdir -p /out/storage/logs /out/database && chmod 777 /out/storage/logs /out/database

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /out/patriot-server /app/patriot-server
COPY --from=build /src/migrations /app/migrations
COPY --from=build /src/configs /app/configs
COPY --from=build --chmod=777 /out/database /app/database
COPY --from=build /out/storage /app/storage
EXPOSE 3000
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 CMD ["/app/patriot-server","-health"]
ENTRYPOINT ["/app/patriot-server"]
