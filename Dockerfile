FROM golang:1.22-alpine AS build

WORKDIR /workspace

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /workspace/api ./cmd/api

FROM alpine:3.20 AS runtime

RUN apk add --no-cache ca-certificates wget \
    && addgroup -S app \
    && adduser -S app -G app

WORKDIR /app

COPY --from=build /workspace/api /app/api

USER app:app

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=5 \
  CMD wget -qO- "http://127.0.0.1:${API_PORT:-8080}/actuator/health/readiness" >/dev/null || exit 1

ENTRYPOINT ["/app/api"]
