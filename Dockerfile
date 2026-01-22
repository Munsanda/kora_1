# ---- build stage ----
FROM golang:1.25.6-alpine AS build

WORKDIR /src

# for fetching modules (and certs if your build downloads over HTTPS)
RUN apk add --no-cache git ca-certificates

# Cache deps first
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source (docs/ is now included since it's not in .dockerignore)
COPY . .

# Build a static binary from cmd/api
ENV CGO_ENABLED=0
RUN go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

# ---- runtime stage ----
FROM alpine:3.22

RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=build /out/api /app/api

# Change if your API listens on another port
EXPOSE 8080

ENTRYPOINT ["/app/api"]