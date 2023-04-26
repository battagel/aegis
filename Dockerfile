FROM golang:1.19

WORKDIR /build

ADD go.mod go.sum Makefile config.env ./
ADD ./internal ./internal
ADD ./pkg ./pkg
ADD ./cmd ./cmd
ADD ./clamd.conf ./clamd.conf

RUN make build

EXPOSE 8080

RUN apt-get update
RUN apt-get install clamdscan -y

RUN cp /build/bin/aegis /usr/local/bin/aegis
CMD ["/usr/local/bin/aegis"]
