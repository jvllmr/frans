# frans

[![Crowdin](https://badges.crowdin.net/go-frans/localized.svg)](https://crowdin.com/project/go-frans)
[![GitHub License](https://img.shields.io/github/license/jvllmr/frans)](https://github.com/jvllmr/frans/blob/main/LICENSE)

Frans is a simple file-sharing service intended to be ready for cloud native.

It took heavy inspiration from [DownloadTicketService](https://www.thregr.org/wavexx/software/dl/). You could also call it a more modern implementation of `DownloadTicketService` optimized for todays (2025) IT landscape.

## Included features

- Create file shares and share them with others
- Create upload grants which allow others to upload their files to you (TODO)

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

Feel free to contribute or improve translations via [!Crowdin](https://crowdin.com/project/go-frans)

## Installation

TODO...

## Requirements

### Database

`frans` supports three types of databases. `PostgreSQL`, `mysql` / `mariadb` and `sqlite`

### Authentication

`frans` does not handle user management by itself, but rather uses OpenID Connect (OIDC) to delegate this task to an OIDC provider. One of the more well known open source examples for an OIDC provider is `Keycloak` which was initially developed by RedHat and has graduated to be a Cloud Native Computing Foundation (CNCF) project. Therefore, `frans` is optimized for usage with `Keycloak`.
If you want to try your luck with another OIDC provider, I have a checklist ready for you. The used OIDC provider needs to:

- provide an introspection endpoint that includes the `sub` claim
- allow usage of refresh tokens
- an `end_session_endpoint`

### SMTP Server

`frans` requires a SMTP server to send mail notifications.

## Configuration

`frans` supports configuration via a `frans.yaml` configuration file relative to its executable or via environment variables.

### Configuration keys

| Key                  | Type     | Description                                                                  | Default                | Environment variable       |
| -------------------- | -------- | ---------------------------------------------------------------------------- | ---------------------- | -------------------------- |
| `root_path`          | string   | The root path from where frans serves its content (i.e. `/files`)            | _empty string_         | `FRANS_ROOT_PATH`          |
| `host`               | string   | IP the webserver should listen on                                            | `127.0.0.1`            | `FRANS_HOST`               |
| `port`               | uint16   | Port the webserver should listen on                                          | `8081`                 | `FRANS_PORT`               |
| `oidc_issuer`        | string   | URL to the OIDC provider                                                     | _(none)_               | `FRANS_OIDC_ISSUER`        |
| `oidc_client_id`     | string   | OIDC client ID                                                               | _(none)_               | `FRANS_OIDC_CLIENT_ID`     |
| `oidc_client_secret` | string   | OIDC client secret                                                           | _(none)_               | `FRANS_OIDC_CLIENT_SECRET` |
| `oidc_admin_group`   | string   | A group in your OIDC provider that should have admin privileges within frans | _empty string_         | `FRANS_DB_PASSWORD`        |
| `db_type`            | string   | Database type to use. One of `postgres`, `mysql` or `sqlite3`                | `postgres`             | `FRANS_DB_TYPE`            |
| `db_host`            | string   | Database host (file path in case of `sqlite3`)                               | `localhost`            | `FRANS_DB_HOST`            |
| `db_port`            | uint16   | Database port                                                                | Database type default  | `FRANS_DB_PORT`            |
| `db_name`            | string   | Database name                                                                | `frans`                | `FRANS_DB_NAME`            |
| `db_user`            | string   | Database user                                                                | `frans`                | `FRANS_DB_USER`            |
| `db_password`        | string   | Database password                                                            | _empty string_         | `FRANS_DB_PASSWORD`        |
| `files_dir`          | string   | Path to directory where files will be stored                                 | `files`                | `FRANS_FILES_DIR`          |
| `max_files`          | uint8    | Max files that can be uploaded per ticket                                    | `20`                   | `FRANS_MAX_FILES`          |
| `max_sizes`          | int64    | Max size per file that can be uploaded                                       | `2_000_000_000` (2 GB) | `FRANS_MAX_SIZES`          |
| `expiry_days_since`  | uint8    | Default expiry days since last download                                      | 7                      | `FRANS_EXPIRY_DAYS_SINCE`  |
| `expiry_total_dl`    | uint8    | Default expiry after total downloads                                         | 10                     | `FRANS_EXPIRY_TOTAL_DL`    |
| `expiry_total_days`  | uint8    | Default expiry after total days                                              | 30                     | `FRANS_EXPIRY_TOTAL_DAYS`  |
| `smtp_server`        | string   | SMTP server host                                                             | _(none)_               | `FRANS_SMTP_SERVER`        |
| `smtp_port`          | integer  | SMTP server port                                                             | `25`                   | `FRANS_SMTP_PORT`          |
| `smtp_username`      | \*string | SMTP username                                                                | _(none)_               | `FRANS_SMTP_USERNAME`      |
| `smtp_password`      | \*string | SMTP password                                                                | _(none)_               | `FRANS_SMTP_PASSWORD`      |
| `smtp_from`          | string   | Default sender email address                                                 | _(none)_               | `FRANS_SMTP_FROM`          |
| `log_json`           | bool     | Whether log messages should be in JSON format                                | `false`                | `FRANS_LOG_JSON`           |

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

#### Build

```shell
pnpm build
```
