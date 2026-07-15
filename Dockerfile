FROM golang:1.25 AS builder

# Install UPX for binary compression
RUN apt-get update && apt-get install -y upx

WORKDIR /go/src/app
COPY . .

RUN go mod download && go vet -v
# Build with optimizations: strip debug info, remove build ID, trim paths
RUN CGO_ENABLED=0 go build -ldflags="-s -w -buildid=false" -trimpath -o /go/bin/app .

# Compress the binary with UPX using ultra-brute for best compression
RUN upx --ultra-brute /go/bin/app

# Use distroless static base for minimal image
FROM gcr.io/distroless/static-debian12

COPY --from=builder /go/bin/app /
EXPOSE 8080
ENTRYPOINT ["/app"]
