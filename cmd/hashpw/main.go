// Package main is a simple utility to generate a bcrypt hash of a password provided via the HASH_PW environment variable.
package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	pw := strings.TrimSpace(os.Getenv("HASH_PW"))
	if pw == "" {
		_, _ = fmt.Fprintln(os.Stderr, "error: HASH_PW environment variable is required")
		os.Exit(1)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// Output base64-encoded hash to avoid $ characters in ldflags
	encoded := base64.StdEncoding.EncodeToString(hash)
	fmt.Print(encoded)
}
