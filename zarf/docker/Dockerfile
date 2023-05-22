FROM golang:1.20-alpine AS builder
ENV CGO_ENABLED 0
ARG BUILD_REF

# Create the service directory and the copy the module files first and then
# download the dependencies. If this doesn't change, we won't need to do this
# again in future builds.
RUN mkdir /cli
COPY go.* /cli/
WORKDIR /cli
RUN go mod download

# Copy the source code into the container.
COPY . /cli

# Build the aws-rds-authenticator binary.
WORKDIR /cli/app/aws-rds-authenticator
RUN go build -ldflags "-X main.build=${BUILD_REF}"

# scratch image
FROM scratch AS scratch

COPY --from=builder /cli/app/aws-rds-authenticator/aws-rds-authenticator .

ENTRYPOINT [ "./aws-rds-authenticator" ]
