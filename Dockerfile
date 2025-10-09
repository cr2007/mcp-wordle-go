FROM golang:1.25 AS builder

WORKDIR /go/src/app
COPY . .

RUN go mod download \
    && go vet -v

RUN CGO_ENABLED=0 go build -o /go/bin/app

FROM gcr.io/distroless/static-debian12

COPY --from=builder /go/bin/app /
CMD ["/app"]
