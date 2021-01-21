FROM golang:alpine

USER root

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN apk add make

WORKDIR /build


COPY src ./
RUN make test && make

CMD ./osoba
