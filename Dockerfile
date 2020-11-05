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
COPY --from=builder bin/ghstsbot /app/
COPY --from=builder entrypoint.sh /app/
RUN chmod +x entrypoint.sh
ENTRYPOINT ["/app/ghstsbot"]
