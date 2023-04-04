FROM aws-rds-authenticator:latest-alpine AS aws-rds-authenticator

FROM alpine:3.17
RUN apk --no-cache add mysql-client

WORKDIR /workspace

COPY bin/entrypoint.sh ./bin/entrypoint.sh
COPY --from=aws-rds-authenticator ./aws-rds-authenticator ./bin/aws-rds-authenticator

ENTRYPOINT [ "./bin/entrypoint.sh" ]
