# Build stage with official Go image
FROM golang:1.20 as builder

WORKDIR /app

# Copy the go mod and sum files
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Runtime stage with UBI Minimal
FROM registry.access.redhat.com/ubi8/ubi-minimal

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/app .

ENTRYPOINT ["./app"]