# Pricing store

## Overview

A simple implementation of a trade pricing store

## Dependencies

- Docker

## Installation for local development

1. `docker compose up` (this builds and deploys up a local Redis container and the pricing store API container)
2. `docker exec pricing-store-redis-1 sh -c "redis-cli < /seed-data/seed.redis"` (this seeds the Redis data store with dummy data, to be automated)

## Endpoints exposed

TBD
