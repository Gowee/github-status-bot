FROM golang:latest as builder
WORKDIR /workspace
COPY ["go.mod", "go.sum", "Makefile", "."]
RUN set -eux; \
    make setup-devdep; \
    make setup-dep
COPY . .
RUN make build

FROM alpine:latest
COPY --from=builder /workspace/bin/ghstsbot /usr/local/bin/ghstsbot
ENTRYPOINT ["/usr/local/bin/ghstsbot"]
