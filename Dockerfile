# Use the official Golang image as a build stage
FROM golang:1.22 as builder

# Set and create and set the Current Working Directory inside the container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download && go mod verify

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o /opt/pipeline

# Copy the Pre-built binary file from the builder stage
FROM alpine:latest
COPY --from=builder /opt/pipeline /opt/pipeline

# Expose port
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["/opt/pipeline"]
