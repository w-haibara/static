FROM golang:alpine AS osoba-builder

USER root
WORKDIR /osoba
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN apk add make

COPY src ./
RUN rm osoba; make test && make

FROM scratch

WORKDIR /osoba

COPY --from=osoba-builder /osoba/osoba .
CMD ["./osoba"]
