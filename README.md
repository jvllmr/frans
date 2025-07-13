# frans

Frans is a simple file-sharing service intended to be ready for cloud native.

It took heavy inspiration from [DownloadTicketService](https://www.thregr.org/wavexx/software/dl/). You could also call it a more modern or v2 version of `DownloadTicketService`.

## Installation

TODO...

## Requirements

### Database

`frans` supports three types of databases. `PostgreSQL`, `mysql` / `mariadb` and `sqlite`

### Authentication

`frans` does not handle user management by itself, but rather uses OpenID Connect (OIDC) to delegate this task to an OIDC provider. One of the more well known open source examples for an OIDC provider is `Keycloak` which was initially developed by RedHat and has to a CNCF project. Therefore, `frans` is optimized for usage with `Keycloak`.
If you want to try your luck with another OIDC provider, I have a checklist ready for you. The used OIDC provider needs to:

- provide an introspection endpoint that includes the `sub` claim
- allow usage of refresh tokens
- an `end_session_endpoint`

### SMTP Server

`frans` requires a SMTP server to send mail notifications.

## Configuration

TODO...

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
