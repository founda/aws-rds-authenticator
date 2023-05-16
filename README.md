<badges>

`aws-rds-authenticator` is a command-line tool that enables users to generate a temporary password for a database
user while leveraging the [AWS RDS IAM Authentication](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.html)
feature. This feature provides an additional layer of security by enabling users to
authenticate with AWS Identity and Access Management (IAM) instead of using a password-based approach. With `aws-rds-authenticator`,
users can easily connect to their database instances using their IAM credentials rather than directly
providing database passwords. This tool simplifies the process and can be particularly useful for Kubernetes users
working with [IAM for Service Accounts](https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts.html).

## Why

The AWS CLI image is big (~125MB) and slow. Pulling and running that image multiple times per day is a waste of resources,
especially for something as trivial as database authentication. Our goal is to make an image under 10MB to fulfil
the same requirement.

## Usage

```bash
$ aws-rds-authenticator -help
Usage of aws-rds-authenticator:
  -database string
      Database that you want to access (optional)
  -engine string
      Database engine that you want to access: postgres|mysql (default "postgres")
  -host string
      Endpoint of the database instance
  -port int
      Port number used for connecting to your DB instance
            default postgres: 5432
            default mysql: 3306
  -region string
      AWS Region where the database instance is running
  -root-cert-file string
      Path to the root certificate file
  -ssl-mode string
      SSL mode to connect to the database instance.
            postgres: disable|require|verify-ca|verify-full (default: verify-ca)
            mysql: DISABLED|PREFERRED|REQUIRED|VERIFY_CA (default: VERIFY_CA)
  -user string
      Database account that you want to access
```

For example, if you want to connect to a PostgreSQL database instance running in the `us-east-1` region, you can use the following command:

```bash
$ aws-rds-authenticator -engine postgres -host rds.amazon.com -port 5432 -user postgres -database prod-db -region us-east-1
<temporary-password>
```

## Authentication

aws-rds-authenticator employs the AWS credentials provider chain by default to authenticate with AWS.

If you're working with Kubernetes, you can use the [IAM for Service Accounts](https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts.html)
feature to authenticate with AWS. EKS will automatically insert the `AWS_WEB_IDENTITY_TOKEN_FILE` and `AWS_ROLE_ARN` environment variables into the pod when the
`eks.amazonaws.com/role-arn` annotation is configured on it.

## Build

```bash
$ go build -o . ./...
```

## Docker

We provide a Docker image for this tool. You can use it like this:

```bash
$ docker build -t aws-rds-authenticator:latest --target=alpine .
```

The following targets are available:

* `alpine`
* `bullseye`
* `scratch`

The default target is `scratch`.

We also provide pre-built images on GitHub Container Registry:

```bash
$ docker run --rm -it ghcr.io/founda/aws-rds-authenticator:latest -help
```

The following tags are available:

* `latest`: Latest stable release, built on scratch
* `latest-alpine`: Latest stable release, built on Alpine Linux
* `latest-bullseye`: Latest stable release, built on Debian Linux

All version tags are available as well (e.g. `1.0.0`, `1.0.0-bullseye`, and `1.0.0-alpine`).

The images are available as linux/amd64 and linux/arm64.

### Using the Docker image

We recommend using the Docker image as a stage in a multi-stage build. This way, you can build your application and copy
the binary to a minimal image.

```dockerfile
FROM ghcr.io/founda/aws-rds-authenticator:latest-alpine AS aws-rds-authenticator
COPY --from=aws-rds-authenticator /workspace/aws-rds-authenticator ./aws-rds-authenticator
```

### Prepared client images

We also provide prepared client images that expose an authenticated database client:

* [founda/aws-rds-authenticator-postgres](https://github.com/founda/aws-rds-authenticator/pkgs/container/aws-rds-authenticator-postgres)
* [founda/aws-rds-authenticator-mysql](https://github.com/founda/aws-rds-authenticator/pkgs/container/aws-rds-authenticator-mysql).

```bash
$ docker pull ghcr.io/founda/aws-rds-authenticator-postgres:latest
$ docker pull ghcr.io/founda/aws-rds-authenticator-mysql:latest
```

These client images require the same variables as the aws-rds-authenticator.

**postgres client entrypoint**

```shell
Usage: ./entrypoint.sh [OPTIONS] [QUERY]

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
    ./entrypoint.sh "SELECT * FROM my_table"
    ./entrypoint.sh -create-db "SELECT * FROM my_table"
```

**mysql client entrypoint**

```shell
Usage: ./entrypoint.sh [OPTIONS] [QUERY]

  Options:
    -create-db  Create the database if it does not exist

  Environment variables:
    MYSQL_HOST      MySQL server host name or IP address
    MYSQL_PORT      MySQL server port number
    MYSQL_USER      MySQL username
    MYSQL_DATABASE  MySQL database name
    AWS_REGION      AWS region where the RDS instance is located

  Examples:
    ./entrypoint.sh "SELECT * FROM my_table"
    ./entrypoint.sh -create-db "SELECT * FROM my_table"
```

### Image size

```bash
ghcr.io/founda/aws-rds-authenticator-mysql      1.0.0            18MB
ghcr.io/founda/aws-rds-authenticator-postgres   1.0.0            10.7MB
ghcr.io/founda/aws-rds-authenticator            1.0.0            5.43MB
ghcr.io/founda/aws-rds-authenticator            1.0.0-alpine     8.69MB
ghcr.io/founda/aws-rds-authenticator            1.0.0-bullseye   51.5MB  

# for reference:
amazon/aws-cli                                  latest           124MB
```

## Alternative

You can also use:

```shell
$ aws rds generate-db-auth-token -hostname my-db-instance.123456789012.us-east-1.rds.amazonaws.com -port 5432 -region us-east-1 -username my-db-username
```

## Contributing

Bug reports, and pull requests are welcome on GitHub. This project is intended to be a safe, welcoming space for
collaboration, and contributors are expected to adhere to the [Contributor Covenant](http://contributor-covenant.org)
code of conduct.

## License

This tool is available as open source under the terms of the [MIT License](http://opensource.org/licenses/MIT).
