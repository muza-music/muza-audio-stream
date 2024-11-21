# Variables
SERVER_BINARY := server
TOKEN_GENERATOR_BINARY := jwtgen
STATIC_DIR := static/
CERTS_DIR := certs
JWT_CERTS_DIR := $(CERTS_DIR)/jwt
SERVER_CERTS_DIR := $(CERTS_DIR)/server
AUDIO_FILES_DIR := audio_files

# Go commands
GO := go
GO_BUILD := $(GO) build
GO_RUN := $(GO) run
GO_MOD := $(GO) mod tidy

# Paths for key files
JWT_PRIVATE_KEY := $(JWT_CERTS_DIR)/private.pem
JWT_PUBLIC_KEY := $(JWT_CERTS_DIR)/public.pem
SERVER_CERT := $(SERVER_CERTS_DIR)/cert.pem
SERVER_KEY := $(SERVER_CERTS_DIR)/key.pem

# Default target
all: build

# Ensure directories exist
$(JWT_CERTS_DIR) $(SERVER_CERTS_DIR) $(AUDIO_FILES_DIR):
	mkdir -p $@

# Generate JWT keys
$(JWT_PRIVATE_KEY) $(JWT_PUBLIC_KEY): | $(JWT_CERTS_DIR)
	openssl genpkey -algorithm RSA -out $(JWT_PRIVATE_KEY) -pkeyopt rsa_keygen_bits:2048
	openssl rsa -pubout -in $(JWT_PRIVATE_KEY) -out $(JWT_PUBLIC_KEY)
	@echo "JWT keys generated: $(JWT_PRIVATE_KEY), $(JWT_PUBLIC_KEY)"

# Generate server SSL certificate and key
$(SERVER_CERT) $(SERVER_KEY): | $(SERVER_CERTS_DIR)
	openssl req -x509 -newkey rsa:2048 -keyout $(SERVER_KEY) -out $(SERVER_CERT) -days 365 -nodes -subj "/CN=localhost"
	@echo "Server SSL certificate and key generated: $(SERVER_CERT), $(SERVER_KEY)"

# Build server and token generator
build: $(SERVER_BINARY) $(TOKEN_GENERATOR_BINARY)

$(SERVER_BINARY): $(JWT_PRIVATE_KEY) $(SERVER_CERT) $(SERVER_KEY)
	$(GO_BUILD) -o $(SERVER_BINARY) cmd/server/main.go
	@echo "Server built: $(SERVER_BINARY)"

$(TOKEN_GENERATOR_BINARY): $(JWT_PRIVATE_KEY)
	$(GO_BUILD) -o $(TOKEN_GENERATOR_BINARY) cmd/jwtgen/main.go
	@echo "Token generator built: $(TOKEN_GENERATOR_BINARY)"

# Run the server
run: build
	./$(SERVER_BINARY)

# Run token generator (example usage)
generate-token: $(TOKEN_GENERATOR_BINARY)
	./$(TOKEN_GENERATOR_BINARY) -username="testuser" -phone="1234" -expiresIn=60

# Tidy up Go modules
tidy:
	$(GO_MOD)

# Clean up generated files
clean:
	rm -f $(SERVER_BINARY) $(TOKEN_GENERATOR_BINARY)
	rm -f $(JWT_PRIVATE_KEY) $(JWT_PUBLIC_KEY)
	rm -f $(SERVER_CERT) $(SERVER_KEY)
	@echo "Cleaned up generated files."

# Help output
help:
	@echo "Usage:"
	@echo "  make              - Build the server and token generator"
	@echo "  make build        - Build both binaries"
	@echo "  make run          - Run the server"
	@echo "  make generate-token - Generate a token with example data"
	@echo "  make tidy         - Tidy up Go modules"
	@echo "  make clean        - Remove built binaries and keys"
	@echo "  make help         - Display this help message"
