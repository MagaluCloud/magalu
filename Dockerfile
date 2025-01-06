# Use Go 1.23 bookworm as base image
FROM golang:1.23-bookworm AS base

# Move to working directory /build
WORKDIR /build

# Copy the entire source code into the container
COPY . .

# Build the application
RUN cd ./mgc/cli && CGO_ENABLED=0 go build -tags "embed release" -o mgc


# Start the application
CMD ["./mgc/cli/mgc", "--version"]

# sudo docker build -t container-registry.br-se1.magalu.cloud/geffteste/mgccli:v0.31.0 . -f Dockerfile
# docker push container-registry.br-se1.magalu.cloud/geffteste/mgccli:v0.31.0