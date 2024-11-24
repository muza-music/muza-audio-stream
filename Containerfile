# Use Red Hat UBI as the base image
FROM registry.access.redhat.com/ubi9/go-toolset:latest AS builder

# Set environment variables
ENV GO111MODULE=on

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files first to leverage Docker cache
COPY go.mod go.sum ./

# Download the Go modules
RUN go mod download

# Copy the entire project
COPY . .

# Change ownership of the workdir to user
USER root
RUN chown -R 1001:1001 /app
USER 1001

# Build the main Go applications with specific flags for minimal environments
RUN go build -o server cmd/server/main.go \
    && go build -o jwtgen cmd/jwtgen/main.go

# Use a smaller UBI image to run the server
FROM quay.io/fedora/fedora:latest

RUN dnf config-manager setopt fedora-cisco-openh264.enabled=1 \
    && dnf install -y https://mirrors.rpmfusion.org/free/fedora/rpmfusion-free-release-$(rpm -E %fedora).noarch.rpm https://mirrors.rpmfusion.org/nonfree/fedora/rpmfusion-nonfree-release-$(rpm -E %fedora).noarch.rpm

RUN sudo dnf swap ffmpeg-freex ffmpeg --allowerasing \
    && dnf clean all

# Set the working directory inside the container
WORKDIR /app

# Copy the built server and jwtgen from the builder image
COPY --from=builder /app/server /app/jwtgen ./

# Copy cert and audio_files directories
COPY certs ./certs
COPY audio_files ./audio_files

# Expose port 8443
EXPOSE 8443

# Run the server as the default command
CMD ["./server"]
