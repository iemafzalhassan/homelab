# Stage 1: Build
FROM golang:1.23-alpine AS builder
ARG SERVICE
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app ./cmd/${SERVICE}/

# Stage 2: Runtime
FROM gcr.io/distroless/static-debian12
COPY --from=builder /app /app
ENTRYPOINT ["/app"]
