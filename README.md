# Amnezia VLESS Builder

[![Build and Release](https://github.com/Lirtich/amnezia-vless-builder/actions/workflows/build-release.yml/badge.svg)](https://github.com/Lirtich/amnezia-vless-builder/actions/workflows/build-release.yml)
[![Tests](https://github.com/Lirtich/amnezia-vless-builder/actions/workflows/tests.yml/badge.svg)](https://github.com/Lirtich/amnezia-vless-builder/actions/workflows/tests.yml)
[![Lint and Format](https://github.com/Lirtich/amnezia-vless-builder/actions/workflows/lint.yml/badge.svg)](https://github.com/Lirtich/amnezia-vless-builder/actions/workflows/lint.yml)

A command-line utility to decrypt and convert Amnezia encrypted provisioning files into usable VLESS connection profiles.

## Overview

Amnezia VLESS Builder is a Go application that processes encrypted VPN configuration files from Amnezia VPN and transforms them into standardized VLESS connection links. It handles the complete pipeline: decryption, decompression, parsing, and profile generation.

### What It Does

1. **Decrypts** provisioning files with the `vpn://` protocol prefix
2. **Decompresses** zlib-compressed configuration data
3. **Parses** the nested JSON configuration structure
4. **Extracts** connection parameters (server, port, user ID, encryption, Reality protocol settings)
5. **Generates** a standardized VLESS URI for use with compatible VPN clients

## Key Features

- 🔐 Supports URL-safe base64 decoding and zlib decompression
- 📋 Extracts all necessary VLESS parameters, including Reality protocol settings
- 🔗 Generates clean VLESS URIs compatible with v2ray, xray, and other clients
- ⚡ Fast and lightweight with minimal dependencies
- 🛠️ Simple command-line interface

## Installation

### Download Pre-built Binaries

Visit the [Releases](https://github.com/Lirtich/amnezia-vless-builder/releases) page to download pre-built binaries for:

- **Linux** (x86_64)
- **macOS** (Intel x86_64 and Apple Silicon ARM64)
- **Windows** (x86_64)

Each release includes compressed archives and SHA256 checksums for verification.

### From Source

#### Requirements

- Go (recent version recommended)

#### Compilation

```
# Build for current platform
go build -o amnezia-vless-builder ./cmd/amnezia-vless-builder
```

Or using `make`:

```
make build          # current platform
make build-all      # all platforms
make build-linux    # Linux x86_64
make build-macos    # macOS (Intel + ARM)
make build-windows  # Windows x86_64
```

Binaries are placed in the `bin/` directory.

## Usage

### Basic Usage

```
./amnezia-vless-builder <source_file>
```

### Example

```
./amnezia-vless-builder provisioning.vpn
```

This will:

- Read the encrypted provisioning file
- Decode and decompress the contents
- Extract all configuration parameters
- Generate a `provisioning.conf` file containing the VLESS connection link
- Print the output path

### Output

The generated `.conf` file contains a VLESS URI in the following format:

```
vless://USER_ID@SERVER:PORT?security=SECURITY&sni=SNI&fp=FINGERPRINT&pbk=PUBLIC_KEY&sid=SHORT_ID&spx=SPIDER_X&type=NETWORK_TYPE&flow=FLOW&encryption=ENCRYPTION#CLIENT_NAME
```

### Parameters Explained

| Parameter | Meaning |
|-----------|---------|
| `USER_ID` | Unique user identifier |
| `SERVER` | VPN server address |
| `PORT` | Connection port |
| `security` | TLS security type (typically "reality") |
| `sni` | Server Name Indication |
| `fp` | TLS fingerprint |
| `pbk` | Reality public key |
| `sid` | Reality short ID |
| `spx` | Spider X path (typically "/") |
| `type` | Network protocol (tcp, ws, etc.) |
| `flow` | VLESS flow control mode |
| `encryption` | Encryption method |
| `CLIENT_NAME` | Profile name (derived from input filename) |

## Project Structure

```
amnezia-vless-builder/
├── cmd/
│   └── amnezia-vless-builder/
│       └── main.go              # Entry point and CLI
├── internal/
│   └── keyconv/
│       ├── decode.go            # Base64 and zlib decoding
│       ├── parse.go             # JSON parsing and parameter extraction
│       ├── process.go           # Main processing pipeline
│       └── render.go            # VLESS URI generation
├── go.mod                        # Go module definition
└── README.md
```

### Module Breakdown

- **main.go**: Command-line interface and argument parsing
- **decode.go**: Handles base64 (URL-safe) and zlib decompression
- **parse.go**: Parses JSON configuration and extracts connection parameters
- **process.go**: Orchestrates the complete file processing workflow
- **render.go**: Generates the final VLESS connection URI

## How It Works

### Processing Pipeline

```
1. Read encrypted file
   ↓
2. Unwrap payload (remove "vpn://" prefix)
   ↓
3. Decode base64 (URL-safe)
   ↓
4. Skip Qt header (4 bytes) and decompress zlib
   ↓
5. Parse outer JSON structure
   ↓
6. Extract inner configuration from xray section
   ↓
7. Extract connection parameters:
   - Server address and port
   - User ID and encryption settings
   - Flow control mode
   - TLS/Reality parameters (SNI, fingerprint, public key, short ID)
   ↓
8. Build VLESS URI
   ↓
9. Save to .conf file
```

### Configuration Structure

The provisioning file contains a nested JSON structure:

```
{
  "containers": [
    {
      "xray": {
        "last_config": "{...inner config JSON...}"
      }
    }
  ]
}
```

The inner configuration contains outbound settings with VLESS protocol parameters and Reality TLS configuration.

## Error Handling

The tool provides clear error messages for common issues:

- **Protocol header absent**: File doesn't start with `vpn://`
- **Base64 decode error**: Invalid base64 encoding
- **Decompression error**: Failed to decompress zlib data
- **Parsing errors**: Invalid JSON structure or missing required fields
- **File errors**: Cannot read source file or write output

## Use Cases

- Convert Amnezia profiles for use with other VLESS clients
- Batch process multiple provisioning files
- Integrate with VPN configuration management systems
- Audit connection parameters from Amnezia profiles

## License

This project is open source. Please check the LICENSE file for details.

## Disclaimer

This tool is provided for educational and legitimate privacy purposes. Users are responsible for complying with local laws and regulations regarding VPN usage in their jurisdiction.

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## Support

For issues, questions, or feature requests, please open an issue on the GitHub repository.

---

**Note**: This tool is designed specifically for Amnezia VPN configuration files. Ensure you have the appropriate provisioning file in the correct format before using this utility.
```
