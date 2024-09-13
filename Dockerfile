
FROM golang:latest as builder

# Set the working directory
WORKDIR /build

ARG GITHUB_ACCESS_TOKEN

# Copy only the Go module files to leverage Docker layer caching
COPY ../go.mod go.sum ./
RUN git config --global url."https://${GITHUB_ACCESS_TOKEN}:@github.com/".insteadOf "https://github.com/" && go mod download 
# Download Go module dependencies

#RUN go mod download

# Copy the entire source code into the container
COPY .. .
COPY cmd/*.go /build
# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main
# Run tests
# RUN go test -v ./...

# Final image, use a lightweight base image
FROM alpine:latest

RUN apk add --no-cache libc6-compat
# Set the working directory
WORKDIR /app

# Copy the built executable from the builder image
COPY --from=builder /app/main .

# Expose the port if your application listens on a specific port
EXPOSE 50052

# Define any environment variables your application may require
# ENV ENV_VAR_NAME=value

# Command to run the application
CMD ["/app/main"]

