# =============================================================================
# Multi-stage build for vault-color-demo
# =============================================================================
FROM golang:1.22-alpine AS builder

WORKDIR /build
COPY go.mod ./
# COPY go.sum ./ # Uncomment when you have dependencies
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o vault-color-demo .

# =============================================================================
# Runtime: distroless for minimal attack surface
# =============================================================================
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app
COPY --from=builder /build/vault-color-demo .
COPY --from=builder /build/templates/ ./templates/

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/vault-color-demo"]
