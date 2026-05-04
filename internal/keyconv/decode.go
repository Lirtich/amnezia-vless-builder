// Package keyconv provides decoding and conversion for Amnezia VPN provisioning files.
package keyconv

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)

// decodeRawBase64 converts URL-safe base64 to binary data.
func decodeRawBase64(enc string) ([]byte, error) {
	trimmed := strings.TrimRight(enc, "=")
	return base64.RawURLEncoding.DecodeString(trimmed)
}

// decompressZlib skips the 4-byte Qt header and decompresses zlib.
func decompressZlib(packed []byte) ([]byte, error) {
	if len(packed) < 4 {
		return nil, fmt.Errorf("packed payload too short")
	}
	zr, err := zlib.NewReader(bytes.NewReader(packed[4:]))
	if err != nil {
		return nil, fmt.Errorf("zlib init: %w", err)
	}
	defer func() {
		_ = zr.Close() //nolint:errcheck // closing a zlib reader is best-effort
	}()
	return io.ReadAll(zr)
}

// unwrapPayload removes "vpn://" prefix, decodes and decompresses the data.
func unwrapPayload(source string) (string, error) {
	trimmed := strings.TrimSpace(source)
	payload, found := strings.CutPrefix(trimmed, "vpn://")
	if !found {
		return "", fmt.Errorf("required protocol header absent")
	}

	raw, err := decodeRawBase64(payload)
	if err != nil {
		return "", fmt.Errorf("base64 decode: %w", err)
	}

	decompressed, err := decompressZlib(raw)
	if err != nil {
		return "", fmt.Errorf("decompression: %w", err)
	}

	return string(decompressed), nil
}
