# frans

Frans is a simple file-sharing service intended to be ready for cloud native.

It took heavy inspiration from [DownloadTicketService](https://www.thregr.org/wavexx/software/dl/). You could also call it a more modern or v2 version of `DownloadTicketService`.

## Authentication

### Requirements

The used OIDC provider needs to:

- provide an introspection endpoint that includes the `sub` claim
- allow usage of refresh tokens
- an `end_session_endpoint`

## Development

### Requirements

- golang installed
- `pnpm` installed
- `docker` with `docker compose v2` installed

### Start environment

#### Start services

```shell
pnpm services
```

This will start the services `frans` depends on for development purposes.
It starts:

- a keycloak instance available under `http://localhost:8080`
- a smtp4dev instance available under `http://localhost:5000`

Keycloak can be managed via the credentials `admin`/`admin`.

Frans authentication is managed via the `dev` realm.

#### Start backend & client

```shell
# Start go backend in dev mode
pnpm dev:go
# Start client dev server
pnpm dev
```

After starting both development servers, `frans` is available under `http://localhost:8081/files`.

You can login via one of the following credentials:

- `frans_admin`/`frans_admin`: User with administration rights
