FROM golang:1.19

WORKDIR /build

ADD go.mod go.sum Makefile config.yaml ./
ADD ./internal ./internal
ADD ./pkg ./pkg
ADD ./vendor ./vendor
ADD ./cmd ./cmd

RUN make build

EXPOSE 8080

RUN cp /build/bin/aegis /usr/local/bin/aegis
CMD ["/usr/local/bin/aegis"]
