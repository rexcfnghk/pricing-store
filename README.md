# Pricing store

## Overview

A simple implementation of a trade pricing store

## Dependencies

- Docker

## Installation for local development

1. `docker compose up` (this builds and deploys up a local Redis container and the pricing store API container)
2. `docker exec pricing-store-redis-1 sh -c "redis-cli < /seed-data/seed.redis"` (this seeds the Redis data store with dummy data, to be automated)

## Endpoints exposed

- [x] `POST /providers/{id}/quotes` Inserts quote prices into the system for a provider
- [x] `GET /providers/{id}/currencyconfigs` Retrieves currency pair setting for a provider
- [ ] `PUT /providers/{id}/currencyconfigs` Updates currency pair setting for a provider

## Things to improve on

- [ ] Lack of tests
- [ ] Code duplication in some areas
- [ ] Lack of sufficient syntax knowledge to construct better abstractions and achieve inversion of control
