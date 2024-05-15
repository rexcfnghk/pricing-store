# syntax=docker/dockerfile:1

FROM golang:latest AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /pricing-store -buildvcs=false
FROM build-stage AS run-test-stage
RUN go test -v ./...
FROM gcr.io/distroless/base-debian11 AS build-release-stage
WORKDIR /
COPY --from=build-stage /pricing-store /pricing-store
COPY --from=build-stage /app/config.json /config.json
COPY --from=build-stage /app/seed-data/ /seed-data/
EXPOSE 8080
CMD ["/pricing-store"]