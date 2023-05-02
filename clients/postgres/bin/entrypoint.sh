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
  exit 1
fi

# Authenticate with AWS RDS and get the PostgreSQL DSN
PG_DSN="$(/workspace/bin/aws-rds-authenticator -engine postgres -host "$PG_HOST" -port "$PG_PORT" -region "$AWS_REGION" -user "$PG_USER")"

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
