<badges>

`aws-rds-authenticator` is a command-line tool that enables users to generate a temporary password for a database user while leveraging the [AWS RDS IAM Authentication](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.html) feature. This feature provides an additional layer of security by enabling users to authenticate with AWS Identity and Access Management (IAM) instead of using a password-based approach. With `aws-rds-authenticator`, users can easily connect to their database instances using their IAM credentials rather than directly providing database passwords. This tool simplifies the process and can be particularly useful for Kubernetes users working with [IAM for Service Accounts](https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts.html).

## Usage

```bash
$ aws-rds-authenticator -help
Usage of aws-rds-authenticator:
  -database string
        Database that you want to access
  -engine string
        Database engine that you want to access: postgres|mysql (default "postgres")
  -host string
        Endpoint of the database instance
  -port int
        Port number used for connecting to your DB instance (default postgres: 5432, default mysql: 3306)
  -region string
        AWS Region where the database instance is running
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

If you're working with Kubernetes, you can use the [IAM for Service Accounts](https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts.html) feature to authenticate with AWS. EKS will automatically insert the `AWS_WEB_IDENTITY_TOKEN_FILE` and `AWS_ROLE_ARN` environment variables into the pod when the `eks.amazonaws.com/role-arn` annotation is configured on it.

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
* `debian`
* `scratch`

The default target is `alpine`.

We also provide pre-built images on GitHub Container Registry:

```bash
$ docker run --rm -it ghcr.io/founda/aws-rds-authenticator:latest -help
```

The following tags are available:

* `latest`: Latest stable release, built on Debian Linux
* `latest-alpine`: Latest stable release, built on Alpine Linux
* `latest-debian`: Latest stable release, built on Debian Linux
* `latest-scratch`: Latest stable release, built on scratch

All version tags are available as well (e.g. `v1.0.0`, `v1.0.0-debian`, and `v1.0.0-alpine`).

### Using the Docker image

We recommend using the Docker image as a stage in a multi-stage build. This way, you can build your application and copy the binary to a minimal image.

```dockerfile
FROM ghcr.io/founda/aws-rds-authenticator:latest-alpine AS aws-rds-authenticator
COPY --from=aws-rds-authenticator /workspace/aws-rds-authenticator ./aws-rds-authenticator
```

## Contributing

Bug reports, and pull requests are welcome on GitHub. This project is intended to be a safe, welcoming space for collaboration, and contributors are expected to adhere to the [Contributor Covenant](http://contributor-covenant.org) code of conduct.

## License

This tool is available as open source under the terms of the [MIT License](http://opensource.org/licenses/MIT).
