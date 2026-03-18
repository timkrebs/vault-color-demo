# =============================================================================
# Multi-stage build for vault-color-demo
# =============================================================================
FROM golang:1.22-alpine AS builder

WORKDIR /build
COPY go.mod ./
# COPY go.sum ./ # Uncomment when you have dependencies
RUN go mod download

COPY . .

# BG_COLOR is injected from Vault via GitHub OIDC at build time
ARG BG_COLOR="#1a1a2e"

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o vault-color-demo .

# =============================================================================
# Runtime: distroless for minimal attack surface
# =============================================================================
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app
COPY --from=builder /build/vault-color-demo .
COPY --from=builder /build/templates/ ./templates/

# Bake the Vault-provided BG_COLOR into the image as default
ENV BG_COLOR=${BG_COLOR}

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/app/vault-color-demo"]
