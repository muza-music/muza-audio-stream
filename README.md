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

Browse to https://localhost:8443 for the demo GAS HTML audio player.

#### Usage

#### Authentication

To authenticate, generate a JWT for the Authorization header in the following format:

Authorization: Bearer <token>

#### Stream Audio

To stream audio, send a POST request to /audio with the following query parameters:

  -  filename: Full path to the audio file you wish to stream.
  -  bitrate: Desired audio bitrate (e.g., 128k).
  -  samplerate: Sample rate of the audio (e.g., 44100).
  -  channels: Number of audio channels (e.g., 2 for stereo).
  -  codec: Audio codec to use (e.g., mp3, aac).
  -  quality: Quality settings for the audio stream (e.g., 4).

#### Example curl Command

```bash
curl -k "https://localhost:8443/audio?filename=audio_files/song.mp3&bitrate=256k&samplerate=48000&channels=2&codec=mp3&quality=4" \
     -H "Authorization: Bearer YOUR_BEARER_TOKEN_HERE"
```

#### Notes

  -  Use -k if you need to bypass SSL certificate validation for local testing (https).
  -  Replace YOUR_BEARER_TOKEN_HERE with your generated JWT.

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

  -  -username: Username to include in the token (required).
  -  -phone: Phone number.
  -  -expiresIn: Number of days until expiration (default: 30)

#### Example:

```bash
# Create a token with expiration date after 60 days (default is 30 days)
./jwtgen -username="john_doe" -phone="054-1234567" -expiresIn=60
```

This will output a JWT token you can use with the streaming server’s Authorization header.
