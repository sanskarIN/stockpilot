package idgen

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func New(prefix string) (string, error) {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", fmt.Errorf("generate id: %w", err)
	}
	if prefix == "" {
		return hex.EncodeToString(raw[:]), nil
	}
	return prefix + "_" + hex.EncodeToString(raw[:]), nil
}
