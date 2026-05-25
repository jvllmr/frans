FROM docker.io/golang:1.26.3-alpine@sha256:91eda9776261207ea25fd06b5b7fed8d397dd2c0a283e77f2ab6e91bfa71079d AS server-base

FROM docker.io/node:24.16.0-alpine@sha256:2bdb65ed1dab192432bc31c95f94155ca5ad7fc1392fb7eb7526ab682fa5bf14 AS client-base
WORKDIR /workspace
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
COPY ./package.json ./pnpm-workspace.yaml ./pnpm-lock.yaml ./
RUN corepack enable pnpm

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
COPY --from=client-builder /workspace/internal ./internal
COPY ./internal/config/version.go ./internal/config/version.go
COPY . .
RUN sh ./scripts/prebuild.sh
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build
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
