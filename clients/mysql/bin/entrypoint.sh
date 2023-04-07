#!/bin/sh

set -e

# Print usage if no arguments are specified
if [ $# -eq 0 ]; then
  echo "Usage: ./entrypoint.sh [OPTIONS] [QUERY]

  Options:
    -create-db  Create the database if it does not exist

  Environment variables:
    MYSQL_HOST      MySQL server host name or IP address
    MYSQL_PORT      MySQL server port number
    MYSQL_USER      MySQL username
    MYSQL_DATABASE  MySQL database name
    AWS_REGION      AWS region where the RDS instance is located

  Examples:
    ./entrypoint.sh \"SELECT * FROM my_table\"
    ./entrypoint.sh -create-db \"SELECT * FROM my_table\""
  exit 0
fi

# Required environment variables
if [ -z "$MYSQL_HOST" ] || [ -z "$MYSQL_PORT" ] || [ -z "$AWS_REGION" ] || [ -z "$MYSQL_USER" ] || [ -z "$MYSQL_DATABASE" ]; then
  echo "ERROR: Missing required environment variables"

  if [ -z "$MYSQL_HOST" ]; then
    echo "  MYSQL_HOST"
  fi

  if [ -z "$MYSQL_PORT" ]; then
    echo "  MYSQL_PORT"
  fi

  if [ -z "$MYSQL_USER" ]; then
    echo "  MYSQL_USER"
  fi

  if [ -z "$MYSQL_DATABASE" ]; then
    echo "  MYSQL_DATABASE"
  fi

  if [ -z "$AWS_REGION" ]; then
    echo "  AWS_REGION"
  fi
  exit 1
fi

# Authenticate with AWS RDS and get the MySQL DSN
MYSQL_DSN="$(/workspace/bin/aws-rds-authenticator -engine mysql -host "$MYSQL_HOST" -port "$MYSQL_PORT" -region "$AWS_REGION" -user "$MYSQL_USER")"

# Create the database if option is specified
if [ "$1" = "-create-db" ]; then
  echo "CREATE DATABASE IF NOT EXIST \`$MYSQL_DATABASE\`;" | mysql "$MYSQL_DSN"
  shift
fi

if [ $# -eq 0 ]; then
  exit 0
fi

echo "$@" | mysql "$MYSQL_DSN" -D "$MYSQL_DATABASE"
