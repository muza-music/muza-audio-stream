# Audio Streaming Server with JWT Authentication

General Audio Server (GAS) provides an audio streaming server using FFmpeg, with JWT-based authentication and support for dynamic audio quality.

## Features
- Supports multiple audio formats (e.g., MP3, WAV, AAC) with FFmpeg transcoding.
- Adjustable audio quality settings.
- Secure streaming with JWT authentication and HTTPS.

## Setup

### 1. Install Dependencies
- **FFmpeg**: Install FFmpeg for audio transcoding.
  
```bash
# Fedora with rpmfusion.org repository enabled, install ffmpeg-free o/w 
sudo dnf install ffmpeg

# Go Modules: Make sure Go modules are enabled:
go mod init gas
go mod tidy
```

### 2. Generate Keys
- Generate JWT Signing Keys

    Private Key: certs/jwt/private.pem
    Public Key: certs/jwt/public.pem

```bash
mkdir -p certs/jwt
openssl genpkey -algorithm RSA -out certs/jwt/private.pem -pkeyopt rsa_keygen_bits:2048
openssl rsa -pubout -in certs/jwt/private.pem -out certs/jwt/public.pem
```

- Generate SSL Certificate and Key for HTTPS

    SSL Key: certs/server/key.pem
    SSL Cert: certs/server/cert.pem

```bash
mkdir -p certs/server
openssl req -x509 -newkey rsa:2048 -keyout certs/server/key.pem -out certs/server/cert.pem -days 365 -
```

### 3. Run the Server

```bash
#Build the Tool
go build -o server cmd/server/main.go

#Wait for streaming requests
./server
```

Browse to https://localhost:8443 for the GAS HTML player.

#### Usage

    Authentication: Generate a JWT for the Authorization header as Bearer <token>.
    Stream Audio: Send requests to /stream with query parameters:
        filenameId: Name of the audio file (without path).
        quality: Audio quality (low, medium, high).

Place audio files in audio_files/, and the server will dynamically transcode and stream them.

## Token Generator Tool

A command-line tool is available to generate JWT tokens with custom usernames.

### Usage

```bash
#Build the Tool
go build -o jwtgen cmd/jwtgen/main.go

#Generate a Token:
./jwtgen -username="your_username"
```

    -username: Username to include in the token (required).

#### Example:

```bash
# Create a token with expiration date after 60 days (default is 30 days)
./jwtgen -username="john_doe" -phone="054-1234567" -expiresIn=60
```

This will output a JWT token you can use with the streaming server’s Authorization header.
