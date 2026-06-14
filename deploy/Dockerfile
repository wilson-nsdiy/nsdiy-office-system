# =============================================================================
# OA-NSDIY Multi-Stage Dockerfile (Compile + Package)
# =============================================================================
# Stage 1: Build frontend (Next.js static export)
# Stage 2: Build Go backend with embedded frontend
# Stage 3: Final minimal runtime image
# =============================================================================

ARG NODE_IMAGE=docker.m.daocloud.io/library/node:20-alpine
ARG GOLANG_IMAGE=docker.m.daocloud.io/library/golang:1.25-alpine
ARG ALPINE_IMAGE=docker.m.daocloud.io/library/alpine:3.20
ARG GOPROXY=https://goproxy.cn,direct
ARG GOSUMDB=sum.golang.google.cn

# -----------------------------------------------------------------------------
# Stage 1: Frontend Builder
# -----------------------------------------------------------------------------
FROM ${NODE_IMAGE} AS frontend-builder

RUN npm config set registry https://registry.npmmirror.com

WORKDIR /app/frontend

COPY frontend/package.json frontend/package-lock.json* ./
RUN npm ci

COPY frontend/ ./
RUN npm run build

# -----------------------------------------------------------------------------
# Stage 2: Backend Builder
# -----------------------------------------------------------------------------
FROM ${GOLANG_IMAGE} AS backend-builder

ARG GOPROXY
ARG GOSUMDB

ENV GOPROXY=${GOPROXY}
ENV GOSUMDB=${GOSUMDB}

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app/backend

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
COPY --from=frontend-builder /app/frontend/out ./internal/web/dist

RUN CGO_ENABLED=0 GOOS=linux go build \
    -tags embed \
    -ldflags="-s -w" \
    -o /app/oa-nsdiy \
    ./cmd/server

# -----------------------------------------------------------------------------
# Stage 3: Final Runtime Image
# -----------------------------------------------------------------------------
FROM ${ALPINE_IMAGE}

LABEL maintainer="OA-NSDIY"
LABEL description="OA-NSDIY - Studio OA Management System"

RUN apk add --no-cache ca-certificates tzdata

RUN addgroup -g 1000 oa-nsdiy && \
    adduser -u 1000 -G oa-nsdiy -s /bin/sh -D oa-nsdiy

WORKDIR /app

COPY --from=backend-builder /app/oa-nsdiy /app/oa-nsdiy

RUN mkdir -p /app/data && chown -R oa-nsdiy:oa-nsdiy /app

EXPOSE 3001

HEALTHCHECK --interval=30s --timeout=10s --start-period=10s --retries=3 \
    CMD wget -q -T 5 -O /dev/null http://localhost:3001/api/health || exit 1

USER oa-nsdiy

CMD ["/app/oa-nsdiy"]
