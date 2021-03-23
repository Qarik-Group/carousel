FROM golang:1.16 as builder

COPY . /code

WORKDIR /code

RUN unset GOPATH && \
    go test -v ./... && \
    CGO_ENABLED=0 go install ./...

FROM alpine

COPY --from=builder /root/go/bin/carousel /bin/

RUN mkdir -p /opt/resource
COPY ./concourse/* /opt/resource/
