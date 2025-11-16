FROM golang:1.25.4-alpine AS server-base

FROM node:24.11.0-alpine AS client-base
WORKDIR /workspace
COPY ./package.json ./pnpm-workspace.yaml ./pnpm-lock.yaml ./
RUN corepack prepare && corepack enable

FROM client-base AS client-deps
WORKDIR /workspace
RUN --mount=type=cache,target=/root/.pnpm-store pnpm install --frozen-lock

FROM client-base AS client-builder
WORKDIR /workspace
COPY ./tsconfig.json ./tsconfig.app.json ./tsconfig.node.json ./vite.config.ts ./
COPY ./locales ./locales
COPY ./client ./client
COPY --from=client-deps /workspace/node_modules ./node_modules
RUN pnpm build:client


FROM server-base AS server-builder
RUN mkdir /emptyd
RUN touch /emptyf
WORKDIR /workspace
ENV CGO_ENABLED=0
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY --from=client-builder /workspace/internal ./internal
COPY ./internal/config/version.go ./internal/config/version.go
COPY . .
RUN sh ./scripts/prebuild.sh
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build go build
RUN chmod +x /workspace/frans


FROM scratch AS runner
COPY --from=server-builder /etc/ssl/certs /etc/ssl/certs
COPY --from=server-builder --chown=1001:1001 /emptyd /tmp
COPY --from=server-builder --chown=1001:1001 /emptyd /opt/frans/files
COPY --from=server-builder --chown=1001:1001 /emptyd /opt/frans/migrations
COPY --from=server-builder --chown=1001:1001 /emptyf /opt/frans/frans.db
COPY --from=server-builder --chown=1001:1001 /emptyf /opt/frans/frans.db-journal
COPY --from=server-builder /workspace/frans /opt/frans/

USER 1001:1001
WORKDIR /opt/frans
VOLUME /opt/frans/files
ENV USER=1001
ENV FRANS_HOST=0.0.0.0
ENV FRANS_DEV_MODE=false

ENTRYPOINT [ "./frans" ]
