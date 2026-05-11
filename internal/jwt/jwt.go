package jwt

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// Decode splits a JWT token into header and payload without verification
func Decode(token string) ([]byte, []byte, error) {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return nil, nil, fmt.Errorf("invalid JWT format")
	}

	header, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, nil, fmt.Errorf("decode header: %w", err)
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, nil, fmt.Errorf("decode payload: %w", err)
	}

	return header, payload, nil
}

// DecodeFull returns header, payload, and signature
func DecodeFull(token string) (map[string]interface{}, map[string]interface{}, string, error) {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return nil, nil, "", fmt.Errorf("invalid JWT format")
	}

	headerJSON, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, nil, "", err
	}

	payloadJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, nil, "", err
	}

	var header, payload map[string]interface{}
	json.Unmarshal(headerJSON, &header)
	json.Unmarshal(payloadJSON, &payload)

	signature := ""
	if len(parts) > 2 {
		signature = parts[2]
	}

	return header, payload, signature, nil
}
