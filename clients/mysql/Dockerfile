ARG VERSION=latest
FROM ghcr.io/founda/aws-rds-authenticator:$VERSION-alpine AS aws-rds-authenticator

FROM alpine:3.21
RUN apk --no-cache add mysql-client

WORKDIR /workspace

COPY bin/entrypoint.sh ./bin/entrypoint.sh
COPY --from=aws-rds-authenticator ./aws-rds-authenticator ./bin/aws-rds-authenticator

ENTRYPOINT [ "./bin/entrypoint.sh" ]
