FROM golang:1.14 as builder
RUN go get -u github.com/pressly/goose/cmd/goose

FROM debian
WORKDIR /app
COPY --from=builder /go/bin/goose .
COPY ./migrations .
RUN chmod 744 migrate.sh
ENTRYPOINT [ "/app/migrate.sh" ]
CMD ["up"]
