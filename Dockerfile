FROM golang:1.16 as builder

COPY . /code

WORKDIR /code

RUN unset GOPATH && \
    go test -v ./... && \
    go install ./...

FROM golang:1.16

RUN mkdir -p /opt/resource

COPY --from=builder /root/go/bin/carousel /bin/
  RUN mkdir -p /opt/resource
  COPY ./concourse/* /opt/resource/
