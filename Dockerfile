# Use the official Golang image as a parent image
FROM golang:1.21.5
# Set the working directory inside the container
WORKDIR /app
# Copy the local package files to the container's workspace.
COPY . .
# Build the Go app
RUN go build -o main ./cmd/app
# Run the binary
CMD ["./main"]