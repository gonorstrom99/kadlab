# FROM alpine:latest

# Add the commands needed to put your compiled go binary in the container and
# run it when the container starts.
#
# See https://docs.docker.com/engine/reference/builder/ for a reference of all
# the commands you can use in this file.
#
# In order to use this file together with the docker-compose.yml file in the
# same directory, you need to ensure the image you build gets the name
# "kadlab", which you do by using the following command:
#
# $ docker build . -t kadlab









# Use the official Golang image as a builder
FROM golang:1.20 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod ./
# COPY go.sum ./

# Download dependencies (this step is cached)
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o kadlab

# Use a smaller image for the final stage to keep the image size down
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/kadlab .

# Expose the port the app runs on (if applicable)
# EXPOSE 8080

# Run the Go app
CMD ["./kadlab"]