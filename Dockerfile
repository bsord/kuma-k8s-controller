# syntax=docker/dockerfile:experimental
FROM golang AS builder
WORKDIR /source
ENV CGO_ENABLED=0
COPY go.* ./
RUN go mod download
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
go build -o=kuma-k8s-controller kuma-k8s-controller

FROM alpine
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

COPY --from=builder "/source/kuma-k8s-controller" /

CMD ["/kuma-k8s-controller"]