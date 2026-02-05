---
title: README
description: A Friendly and Selfhostable homepage, Tree style bookmarks, Widgets for popular selfhosted sites, and more to come!
author: Adam Kalinowski
created: 2026-01-05T14:48:44-0500
updated: 2026-01-05T15:07:41-0500
version: 1.1.1
---


# Mindscape: A Friendly And Selfhostable homepage

**Mindscape** is a selfhostable homepage built specifically for those who loved the ARC broweser before it was deprecated in favor of some ai slop browser. I loved teh ARC browser because of three things:
- There was a structure to bookmarks and links in the way that programmers generally think about information: a **tree**
- There was a clear seperation between workspaces (work, opensource, and content creation) frequently used sites
- There was a clear focus on developing with public apis (like pull requests on your github account) that was hopefully going to mature into a ecosystem that people would love to use and contribute to.
  But like sillicon valley startups the dream was not meant to be and the dragon was chased and we were left with nothing as an alternative.

## So why this


## How to install

### Install

It is recomended that you use docker compose to host the mindscape server. Here are some options you can use.

#### Use the following docker compose file:

```yaml
services:
  web:
    image: ghcr.com/adamkali/mindscape
    container_name: mindscape
    ports:
      - "60000:60000"
    volumes:
      - app:/app
    healthcheck:
      test:
        - CMD
        - curl
        - -f
        - http://127.0.0.1:60000/_health
      interval: 5s
      timeout: 20s
      retries: 10
  db:
    image: postgres
    container_name: mindscape-db
orts:
 - "60001:5432"
nvironment:
      POSTGRES_PASSWORD: <PASSWORD> 
 POSTGRES_USER: <USERNAME>
      POSTGRES_DB: mindscape
ealthcheck:
 test:
   - CMD
pg_isready
 interval: 10s
      timeout: 5s
      retries: 10

  redis:
    image: redis # or any redis equivalent
    container_name: mindscape-cache
orts:
 - "60002:6379"
    environment:
      REDIS_PASSWORD: <PASSWORD>
      REDIS_USER: <USERNAME>
ealthcheck:
     test:
        - CMD
        - redis-cli
   - ping
 interval: 10s
 timeout: 5s
      retries: 10
  s3:
mage: minio/minio
ontainer_name: mindscape-s3
orts:
 - "60003:9000"
nvironment:
 MINIO_ROOT_USER: <USERNAME>
 MINIO_ROOT_PASSWORD: <PASSWORD>
ealthcheck:
 test:
   - CMD
   - curl
   - -f
        - http://127.0.0.1:60003/minio/health/live
```
  And then to run it, use:
```bash
docker compose up -d 
```
  You will also have to source the config file, associated in [Config](#config) below and point the database an such to the correct location.

#### Coolify

> TODO: Need to add instructions for coolify

### Config

In order to get Mindscape to work, yo will need to specify where your database lives, where the cache is living and where s3 is living,
```yaml
namespace: github.com/adamkali/mindscape # <-- you can leave this as it is and should not affect runtime
name: mindscape
semver: 0.0.5
license: MIT
copyright:
  year: 2025
  author: Adam Kali
server:
wt: # 
  port: 60000
  frontend:
    dir: web/dist
    api: web/src/api
database:
rl: # 
  sqlc:
    repository: db/repository
    schema: postgresql
    sql_or_go: sql
  queries: db/queries
  migration:
    protocol: postgres
    destination: db/migrations
cache:
rl: redis://:#@snickers:60002/0
s3:
  url: minio-j8s8okssws040w0goggko4gg.kalilarosa.xyz
ccess: # 
ecret: # 
```
