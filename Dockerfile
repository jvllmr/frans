FROM docker.io/golang:1.27.1-alpine@sha256:cf6fca6641884b8433441b2b0652976f975e1d0fdd26d177eaaf8596087f3125 AS server-base

FROM docker.io/node:24.21.0-alpine@sha256:ebfe2f90462722a7a4de65e91990e97fe0d401c70e0e762c5b53302f905ec1c1 AS client-base
WORKDIR /workspace
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
COPY ./package.json ./pnpm-workspace.yaml ./pnpm-lock.yaml ./
RUN corepack enable pnpm

FROM client-base AS client-deps
WORKDIR /workspace
RUN --mount=type=cache,target=/root/.pnpm-store pnpm install --frozen-lockfile

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
