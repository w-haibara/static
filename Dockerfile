FROM golang:alpine AS builder

USER root
WORKDIR /osoba
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN apk add make

COPY src ./
RUN make test && make

FROM scratch

WORKDIR /osoba

COPY --from=builder /osoba/osoba .
CMD ["./osoba"]
