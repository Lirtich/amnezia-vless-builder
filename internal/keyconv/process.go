// Package keyconv provides decoding and conversion for Amnezia VPN provisioning files.
package keyconv

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ProcessFile performs the full transformation of the source file and writes the result.
func ProcessFile(inputPath string) (string, error) {
	raw, err := os.ReadFile(inputPath)
	if err != nil {
		return "", fmt.Errorf("cannot read source: %w", err)
	}

	payload, err := unwrapPayload(string(raw))
	if err != nil {
		return "", fmt.Errorf("unpacking: %w", err)
	}

	inner, err := parseInnerConfig(payload)
	if err != nil {
		return "", fmt.Errorf("parsing config: %w", err)
	}

	params, err := extractParams(inner)
	if err != nil {
		return "", fmt.Errorf("extracting fields: %w", err)
	}

	clientName := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	link := buildProfileLink(clientName, params)

	outPath := strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + ".conf"
	if err := os.WriteFile(outPath, []byte(link), 0644); err != nil {
		return "", fmt.Errorf("writing profile: %w", err)
	}

	return outPath, nil
}
