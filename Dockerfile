FROM golang AS builder
WORKDIR /source
COPY . .
RUN go get ./...
RUN go build -v -o=kuma-k8s-operator kuma-k8s-operator/client
RUN mkdir /app && cp kuma-k8s-operator /app/

FROM alpine
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
WORKDIR /app

COPY --from=builder "/app" .

CMD ["/app/kuma-k8s-operator", "run"]