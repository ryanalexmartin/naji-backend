# Use the official Go image as a base image
FROM golang:1.17

# Set the working directory
WORKDIR /app

# Copy the Go modules and build files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the application
RUN go build -o main

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["/app/main"]
