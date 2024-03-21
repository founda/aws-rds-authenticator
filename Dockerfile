# syntax=docker/dockerfile:1.4
FROM golang:1.20-alpine AS builder

WORKDIR /workspace

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Build the binary
RUN CGO_ENABLED=0 go build -o . ./...

# scratch image
FROM scratch AS scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /workspace/aws-rds-authenticator .

ENTRYPOINT [ "./aws-rds-authenticator" ]

# debian image
FROM debian:bullseye-slim AS bullseye

RUN apt-get update && apt-get install -y ca-certificates
RUN addgroup --system app --gid 888 && \
    adduser --system --no-create-home --uid 888 --ingroup app app

COPY --from=builder --chown=app:app /workspace/aws-rds-authenticator .

USER app

ENTRYPOINT [ "./aws-rds-authenticator" ]

# alpine image
FROM alpine:3.17 AS alpine

RUN addgroup --system app --gid 888 && \
    adduser --system --no-create-home --uid 888 --ingroup app app

COPY --from=builder --chown=app:app /workspace/aws-rds-authenticator .

USER app

ENTRYPOINT [ "./aws-rds-authenticator" ]
