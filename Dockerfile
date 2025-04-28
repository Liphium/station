# Use official Golang image as base image
FROM golang:latest as builder
WORKDIR /app
COPY . .
RUN go mod download
WORKDIR /app/main

# Build the Go app
RUN go build -o start .

# Make the final container based on a small image
FROM debian:bookworm
WORKDIR /app

# Make sure certificates work properly (for R2 and decentralization)
RUN apt-get update
RUN apt-get install -y ca-certificates
RUN update-ca-certificates

# Copy the current executable over to the container from the builder
COPY --from=builder /app/main/start .

# Run the app together with the ports
EXPOSE 3000 3001 3002
CMD ["./start"]
