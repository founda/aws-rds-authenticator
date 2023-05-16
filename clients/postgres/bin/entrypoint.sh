#!/bin/sh

set -e

# Print usage if no arguments are specified
if [ $# -eq 0 ]; then
  echo "Usage: ./entrypoint.sh [OPTIONS] [QUERY]

  Options:
    -create-db  Create the database if it does not exist

  Environment variables:
    PG_HOST     PostgreSQL server host name or IP address
    PG_PORT     PostgreSQL server port number
    PG_USER     PostgreSQL username
    PG_DATABASE PostgreSQL database name
    AWS_REGION  AWS region where the RDS instance is located
    PG_SSL_MODE (optional) PostgreSQL SSL mode (default: verify-ca)
    PG_SSL_CA   (optional) PostgreSQL SSL CA certificate file path (required if PG_SSL_MODE is verify-ca or verify-full)

  Examples:
    ./entrypoint.sh \"SELECT * FROM my_table\"
    ./entrypoint.sh -create-db \"SELECT * FROM my_table\""
  exit 0
fi

# Required environment variables
if [ -z "$PG_HOST" ] || [ -z "$PG_PORT" ] || [ -z "$AWS_REGION" ] || [ -z "$PG_USER" ] || [ -z "$PG_DATABASE" ]; then
  echo "ERROR: Missing required environment variables"

  if [ -z "$PG_HOST" ]; then
    echo "  PG_HOST"
  fi

  if [ -z "$PG_PORT" ]; then
    echo "  PG_PORT"
  fi

  if [ -z "$PG_USER" ]; then
    echo "  PG_USER"
  fi

  if [ -z "$PG_DATABASE" ]; then
    echo "  PG_DATABASE"
  fi

  if [ -z "$AWS_REGION" ]; then
    echo "  AWS_REGION"
  fi

  # if PG_SSL_MODE is not set, default to verify-ca
  if [ -z "$PG_SSL_MODE" ]; then
    PG_SSL_MODE="verify-ca"
  fi

  # if PG_SSL_MODE is verify-ca or verify-full, PG_SSL_CA is required
  if [ "$PG_SSL_MODE" = "verify-ca" ] || [ "$PG_SSL_MODE" = "verify-full" ]; then
    if [ -z "$PG_SSL_CA" ]; then
      echo "  PG_SSL_CA"
    fi
  fi

  exit 1
fi

# Authenticate with AWS RDS and get the PostgreSQL DSN
ARGS="-engine postgres -host $PG_HOST -port $PG_PORT -region $AWS_REGION -user $PG_USER"
if [ -n "$PG_SSL_MODE" ]; then
  ARGS="$ARGS -ssl-mode $PG_SSL_MODE"
fi
if [ -n "$PG_SSL_CA" ]; then
  ARGS="$ARGS -root-cert-file $PG_SSL_CA"
fi

PG_DSN="$(/workspace/bin/aws-rds-authenticator "$ARGS")"

# Create the database if option is specified
if [ "$1" = "-create-db" ]; then
  # see: https://stackoverflow.com/a/18389184
  echo "SELECT 'CREATE DATABASE $PG_DATABASE' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$PG_DATABASE')\gexec" | psql "$PG_DSN"
  shift
fi

if [ $# -eq 0 ]; then
  exit 0
fi

echo "$@" | psql "$PG_DSN" -d "$PG_DATABASE"
