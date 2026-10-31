FROM docker.io/golang:1.27.3-alpine@sha256:d200d2dd76b03a99b52651fd1325aff7dcd35422c35c01bac56571b0507f3deb AS server-base

FROM docker.io/node:24.22.2-alpine@sha256:8a32c8cf649a49915a66e81e7a609223c708d0adabf288ee3f4077d765542f85 AS client-base
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
WORKDIR /workspace
ENV CGO_ENABLED=0
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY ./internal/config/version.go ./internal/config/version.go
COPY . .
COPY --from=client-builder /workspace/internal/routes/client/assets ./internal/routes/client/assets
RUN sh ./scripts/prebuild.sh
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build
RUN mkdir -p /rootfs/opt/frans/files /rootfs/opt/frans/migrations /rootfs/tmp \
    && cp /workspace/frans /rootfs/opt/frans/frans \
    && chmod +x /rootfs/opt/frans/frans \
    && touch /rootfs/opt/frans/frans.db /rootfs/opt/frans/frans.db-journal


FROM scratch AS runner
COPY --from=server-builder /etc/ssl/certs /etc/ssl/certs
COPY --from=server-builder --chown=1001:1001 /rootfs/ /

USER 1001:1001
WORKDIR /opt/frans
VOLUME /opt/frans/files
ENV USER=1001
ENV FRANS_HOST=0.0.0.0
ENV FRANS_DEV_MODE=false

ENTRYPOINT [ "./frans" ]
