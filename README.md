# Pricing store ![CI Build](https://github.com/rexcfnghk/pricing-store/actions/workflows/go.yml/badge.svg)

## Overview

A simple implementation of a trade pricing store

## Dependencies

- Docker

## Installation for local development

1. `docker compose up`
   - This builds and deploys
      - A Redis container hosted on port `6379`
      - The pricing API hosted on port `3000`
      - A Swagger UI hosted on port `80`

## Endpoints exposed

- [x] `POST /providers/{id}/quotes` Inserts quote prices into the system for a provider
- [x] `GET /providers/{id}/currencyconfigs?base={baseCurrency}&quote={quoteCurrency}` Retrieves currency pair setting for a provider
- [x] `PUT /providers/{id}/currencyconfigs?base={baseCurrency}&quote={quoteCurrency}` Updates currency pair setting for a provider
- [x] `GET /providers/bestprice?base={baseCurrency}&quote={quoteCurrency}` Get best price among all available providers for a given currency pair
  - Requires a JWT with a payload `{ "sub": {customerId} }`

The endpoints can also be visualised by visiting the Swagger UI generated from `openapi/openapi.yaml`

## Assumptions

- Provider IDs are known to the caller and are provided as URL parameters
- Customer IDs are known to the caller and are provided as JWTs

## Seed data

`docker compose up` defaults to seeding the Redis instance with data located under `seed-data/seed.redis`. Seed data can be modified by updating the contents of this file.

## Things to improve on

- [x] Add Postman scripts for easier testing
- [ ] Lack of tests
  - Due to a lack of time, only one simple unit test is added under `repository/customer/redis_test.go` to demonstrate a unit test written in Go
- [ ] Code duplication in some areas
- [ ] Lack of sufficient syntax knowledge to construct better abstractions and achieve inversion of control
