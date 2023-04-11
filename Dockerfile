# syntax=docker/dockerfile:1

FROM golang:1.20 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux go build -o /storage

FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /storage /storage

USER nonroot:nonroot

ENTRYPOINT ["/storage"]