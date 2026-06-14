FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X main.Version=${VERSION}" \
    -trimpath \
    -o server ./cmd/server

FROM alpine:3.21

RUN apk --no-cache add ca-certificates tzdata

RUN addgroup -g 1000 app && \
    adduser -u 1000 -G app -s /bin/sh -D app

WORKDIR /app

COPY --from=builder --chown=app:app /app/server .
COPY --from=builder --chown=app:app /app/config.example.yaml .

RUN mkdir -p /app/data && chown app:app /app/data

EXPOSE 3001

HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
    CMD wget -q -T 5 -O /dev/null http://localhost:3001/api/health || exit 1

USER app

CMD ["./server"]
