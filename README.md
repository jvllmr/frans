# frans

[![Crowdin](https://badges.crowdin.net/go-frans/localized.svg)](https://crowdin.com/project/go-frans)
[![GitHub License](https://img.shields.io/github/license/jvllmr/frans)](https://github.com/jvllmr/frans/blob/main/LICENSE)

Frans is a simple file-sharing service intended to be ready for cloud native.

It took heavy inspiration from [DownloadTicketService](https://www.thregr.org/wavexx/software/dl/). You could also call it a more modern implementation of `DownloadTicketService` optimized for todays (2025) IT landscape.

## Included features

- Create file shares and share them with others
- Create upload grants which allow others to upload their files to you

The goal is to translate `frans` to the following languages:

- `German`
- `French`
- `Spanish`
- `Russian`
- `Simplified chinese`
- `Japanese`
- `Dutch`
- `Portuguese, Brazilian`
- `Czech`
- `Italian`

Feel free to contribute or improve translations via [Crowdin](https://crowdin.com/project/go-frans)

## Installation

### Helm Chart

TODO...

### Docker Compose

See [docker-compose.minimal.yaml](docker-compose.minimal.yaml) for a minimal example.

### Pre-compiled binaries

Take a look at [Frans GitHub Releases](https://github.com/jvllmr/frans/releases) for pre-compiled binaries for linux, macOS and windows!

## Requirements

### Database

`frans` supports three types of databases. `PostgreSQL`, `mysql` / `mariadb` and `sqlite`

### Authentication

`frans` does not handle user management by itself, but rather uses OpenID Connect (OIDC) to delegate this task to an OIDC provider. One of the more well known open source examples for an OIDC provider is `Keycloak` which was initially developed by RedHat and has graduated to be a Cloud Native Computing Foundation (CNCF) project. Therefore, `frans` is optimized for usage with `Keycloak`.
If you want to try your luck with another OIDC provider, I have a checklist ready for you. The used OIDC provider client needs to:

- be a public client (no client secret needed) with pkce challenges
- allow usage of refresh tokens
- have an `end_session_endpoint`

### SMTP Server

`frans` requires a SMTP server to send mail notifications.

## Configuration

`frans` supports configuration via a `frans.yaml` configuration file relative to its executable or via environment variables.

### Configuration keys

See [frans.defaults.yaml](frans.defaults.yaml)

## Development

I always welcome contributions. Even the smaller ones.

### Development requirements

The following tools are required to start with development on frans:

- golang installed
- `pnpm` installed: installation via corepack is preferred
- `docker` with `docker compose v2` installed
- `atlas` installed: <https://atlasgo.io/guides/evaluation/install>

### Start environment

After you have installed all requirements you can start the environment

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

#### Setup database

```shell
pnpm db:migrate
```

#### Start backend & client for development

```shell
# Start go backend in dev mode
pnpm dev:go
# Start client dev server
pnpm dev
```

After starting both development servers, `frans` is available under `http://localhost:8081/files`.

You can login via one of the following credentials:

- `frans_admin`/`frans_admin`: User with administration rights

### Do changes to database

Everything for ent is generated from the contents of `internal/ent/schema`.
After making your changes there, run the following to apply them correctly:

```shell
# generate new ent source
pnpm ent:generate

# create migration scripts for postgresql, mysql and sqlite3
pnpm db:diff

# apply migrations to development database
pnpm db:migrate
```

#### Create a new database entity

```shell
pnpm ent:new <entity name>
```

#### Build

```shell
pnpm build
```
