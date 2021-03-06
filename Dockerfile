FROM golang:latest as builder
WORKDIR /workspace
COPY ["go.mod", "go.sum", "Makefile", "./"]
RUN set -eux; \
    make setup-devdep; \
    make setup-dep
COPY . .
RUN make build

FROM alpine:latest
WORKDIR /app
VOLUME ["/app/data"]
COPY --from=builder /workspace/bin/ghstsbot /app/ghstsbot
COPY --from=builder /workspace/entrypoint.sh /app/entrypoint.sh
RUN chmod +x entrypoint.sh
ENTRYPOINT ["sh", "entrypoint.sh"]
