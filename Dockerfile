FROM golang:1.25 AS builder

# Install UPX for binary compression, and CA certs for outbound HTTPS
# (the app calls the NYT Wordle API), since the final image has no package manager to fetch them.
RUN apt-get update && apt-get install -y upx ca-certificates

WORKDIR /go/src/app
COPY . .

RUN go mod download && go vet -v
# Build with optimizations: strip debug info, remove build ID, trim paths
RUN CGO_ENABLED=0 go build -ldflags="-s -w -buildid=false" -trimpath -o /go/bin/app .

# Compress the binary with UPX using ultra-brute for best compression
RUN upx --ultra-brute /go/bin/app

# scratch has no tzdata/netbase/passwd/etc, which distroless bundles but this
# app never uses (no time zone handling, single static binary, no users/groups).
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/bin/app /
EXPOSE 8080
ENTRYPOINT ["/app"]
