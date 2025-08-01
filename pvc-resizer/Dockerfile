# Stage 1: Build Go binary
FROM golang:1.24 AS builder
WORKDIR /app

COPY . .
# 👇 This ensures the binary works in the Linux container
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o pvc-operator main.go

# Stage 2: Use distroless image
FROM gcr.io/distroless/static:latest
COPY --from=builder /app/pvc-operator /pvc-operator
ENTRYPOINT ["/pvc-operator"]
