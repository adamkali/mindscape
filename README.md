---
title: README
description: A friendly, self-hostable homepage — tree-style bookmarks, notes, widgets for popular self-hosted sites, and a keyboard-driven mindpalace for your internet.
author: Adam Kalinowski
created: 2026-01-05T14:48:44-0500
updated: 2026-06-05
version: 1.2.0
---

# Mindscape — a self-hostable case board for your internet

![Mindscape hero](assets/launch/hero.png)

<video src="assets/launch/demo.mp4" controls muted loop width="100%">
  Your browser does not support the video tag — <a href="assets/launch/demo.mp4">watch the demo</a>.
</video>

**Mindscape** is a self-hostable homepage built for people who think about their
internet as a *structure*, not a flat pile of tabs. Tree-style bookmarks, notes,
and research live in one keyboard-driven space you don't just look at — you
*interact* with it.

There are so many great pages out there, I dont want to have to search for the 
specific one i want every single time. I know what I want to go to. And bookmarks in browsers
always seemed like a little lacking in my opinion.

If this sounded like you. Then you most likely loved the arc browser. 
It is built for those who loved the ARC browser before it was deprecated in favor
of some AI-slop browser. I loved the ARC browser because of three things:

- There was a structure to bookmarks and links in the way that programmers
  generally think about information: a **tree**.
- There was a clear separation between workspaces (work, open source, and content
  creation) and frequently used sites.
- There was a clear focus on developing with public APIs (like pull requests on
  your GitHub account) that was hopefully going to mature into an ecosystem that
  people would love to use and contribute to.

The dream was not meant to be, the dragon was chased, and we were left with nothing
as an alternative.

## So why this

The real inspiration isn't a browser — it's **Saga Anderson's Mind Place / case
board from Alan Wake 2**: a detective's evidence board where you pin what matters
and walk the connections between things. Mindscape is that idea for your actual
internet: a **mindpalace** you build out of folders, bookmarks, and notes, then
navigate entirely from the keyboard.

The path was *case board → mindpalace → Mindscape*. Top-level folders become your
workspaces; nested folders branch into projects; bookmarks and markdown notes are
the evidence you pin to each branch. Everything is self-hosted, so the board is
yours — no cloud account, no telemetry, no rug-pull. On top of this you don't 
need to move your book marks around in the event of changing to a new browser.
Point to your mindscape instance as a new tab page, and then you are ready.


## How to install

It is recommended that you use Docker Compose to host the Mindscape server. Here
are some options you can use.

### Use the following docker compose file

```yaml
services:
  web:
    image: ghcr.io/adamkali/mindscape
    container_name: mindscape
    ports:
      - "60000:60000"
    volumes:
      - app:/app/data
      # Supply your own production config (see the Config section below).
      - ./config/production.yaml:/app/config/production.yaml:ro
    depends_on:
      db:
        condition: service_healthy
      redis:
        condition: service_healthy
      s3:
        condition: service_healthy
    healthcheck:
      test:
        - CMD
        - curl
        - -f
        - http://127.0.0.1:60000/api/_health
      interval: 5s
      timeout: 20s
      retries: 10

  db:
    image: postgres:16
    container_name: mindscape-db
    ports:
      - "60001:5432"
    environment:
      POSTGRES_PASSWORD: <PASSWORD>
      POSTGRES_USER: <USERNAME>
      POSTGRES_DB: mindscape
    volumes:
      - db:/var/lib/postgresql/data
    healthcheck:
      test:
        - CMD
        - pg_isready
        - -U
        - <USERNAME>
      interval: 10s
      timeout: 5s
      retries: 10

  redis:
    image: redis:7 # or any redis equivalent
    container_name: mindscape-cache
    ports:
      - "60002:6379"
    command:
      - redis-server
      - --requirepass
      - <PASSWORD>
    healthcheck:
      test:
        - CMD
        - redis-cli
        - ping
      interval: 10s
      timeout: 5s
      retries: 10

  s3:
    image: minio/minio
    container_name: mindscape-s3
    command:
      - server
      - /data
      - --console-address
      - ":9001"
    ports:
      - "60003:9000"
      - "60004:9001"
    environment:
      MINIO_ROOT_USER: <USERNAME>
      MINIO_ROOT_PASSWORD: <PASSWORD>
    volumes:
      - s3:/data
    healthcheck:
      test:
        - CMD
        - curl
        - -f
        - http://127.0.0.1:9000/minio/health/live
      interval: 10s
      timeout: 5s
      retries: 10

volumes:
  app: {}
  db: {}
  redis: {}
  s3: {}
```

And then, to run it, use:

```bash
docker compose up -d
```

You will also have to source the config file described in [Config](#config) below
and point the database, cache, and S3 URLs at the correct locations. Run the
database migrations once on first boot:

```bash
docker compose exec web ./mindscape migrate up
```

#### Coolify

> TODO: Coolify deploy instructions are tracked separately and will land in a
> follow-up.

### Config

In order to get Mindscape to work, you will need to specify where your database
lives, where the cache is living, and where S3 is living. The keys below map
one-to-one to the server's configuration struct. When running under the Compose
file above, address the dependencies by their service name and internal port.

```yaml
namespace: github.com/adamkali/mindscape # you can leave this as-is; it does not affect runtime
name: mindscape
semver: 0.0.3
license: MIT
copyright:
  year: 2025
  author: Adam Kali
server:
  jwt: <a-long-random-secret> # used to sign JWTs — generate your own
  port: 60000
  frontend:
    dir: web/dist
    api: web/src/api
database:
  url: postgres://<USERNAME>:<PASSWORD>@db:5432/mindscape
  sqlc:
    repository: db/repository
    schema: postgresql
    sql_or_go: sql
  queries: db/queries
  migration:
    protocol: postgres
    destination: db/migrations
cache:
  url: redis://:<PASSWORD>@redis:6379/0
s3:
  url: s3:9000
  access: <USERNAME>
  secret: <PASSWORD>
  secure: false
apikey:
  default_expiration: 2592000000000000
households:
  invite_ttl_hours: 168
features:
  tasks: false
  apikeys: false
  coolify: false
```

### Feature flags

Some surfaces ship **off** in the v1 build. They are flagged in two places,
because the server and the browser bundle get their configuration by different
routes — the server reads YAML at startup, the frontend inlines values at build
time. **Both sides must agree**, or you get a button that 404s.

| Surface | Server (`config/*.yaml`) | Frontend (build-time env) | Default |
|---|---|---|---|
| Tasks / Agenda | `features.tasks` | `MINDSCAPE_FEATURE_TASKS=true` | off |
| API keys | `features.apikeys` | `MINDSCAPE_FEATURE_APIKEYS=true` | off |
| Coolify widget | `features.coolify` | *(server only)* | off |
| Admin pages | *(n/a)* | `MINDSCAPE_ADMIN=true` | on in dev, off in prod builds |

Notes on the flags:

- **Server side is YAML only.** There is no env-override layer in
  `cmd/configuration/configuration.go`, so `features.tasks` cannot be set by an
  environment variable — edit the config file you mount into the container.
- **Turning `features.apikeys` off disables the whole `/api/key/*` group**, not
  just the API-key CRUD endpoints. Key auth is the only way to authenticate to
  that group, so its folder and bookmark routes would be unreachable anyway.
- **The Coolify widget stays in the binary either way.** Its schema is embedded
  at compile time and simply held out of the default install list; setting
  `features.coolify: true` puts it back in the Add Widget picker.
- Frontend flags are read in `web/src/config/features.ts` and injected by
  `source.define` in `web/rsbuild.config.ts`. Rsbuild has no Vite-style
  `import.meta.env`, so add any new flag in both files.

## Links

- **Blog post:** [blog.kalilarosa.xyz](https://blog.kalilarosa.xyz)
- **Issues / feedback:** [github.com/adamkali/mindscape/issues](https://github.com/adamkali/mindscape/issues)
- **Launch demo:** [watch the 60-second tour](assets/launch/demo.mp4)

## License

Mindscape is released under the [MIT License](LICENSE).
