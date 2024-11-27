# # syntax=docker/dockerfile:1

# FROM golang:1.20

# # Set destination for COPY
# WORKDIR /goSync

# # Enviroment Application Properties
# COPY .env ./
# # Download Go modules
# COPY go.mod go.sum ./
# RUN go mod tidy
# RUN go mod download

# # Copy the rest of the code
# COPY . .

# # Build
# RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-go-sync

# # Optional:
# # To bind to a TCP port, runtime parameters must be supplied to the docker command.
# # But we can document in the Dockerfile what ports
# # the application is going to listen on by default.
# # https://docs.docker.com/reference/dockerfile/#expose
# #EXPOSE 8080

# # Run
# CMD ["/docker-go-sync"]

# Step 1: Build the Go application using the official Go image
FROM golang:1.20 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to leverage Docker's caching for dependencies
COPY go.mod go.sum ./

COPY .env ./

# Download all Go module dependencies
RUN go mod download


# Copy the rest of the application code
COPY . .

# Build the Go application (CGO disabled for a fully static binary, targeting Linux)
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-app

# # Step 2: Use a minimal base image for the final container (e.g., Alpine)
# FROM alpine:latest

# # Set a working directory inside the container (optional)
# WORKDIR /root/

# # Copy the Go binary from the builder container
# COPY --from=builder /go-app .

# Set the default command to run the binary
CMD ["./go-app"]