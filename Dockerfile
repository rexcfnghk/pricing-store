# syntax=docker/dockerfile:1

FROM golang:latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /pricing-store
EXPOSE 8080
CMD ["/pricing-store"]