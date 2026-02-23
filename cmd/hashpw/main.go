// Package main is a simple utility to generate a bcrypt hash of a password provided via the HASH_PW environment variable.
package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var Srvc = ""

func main() {
	switch Srvc {
	case "scale":
		hashPassword("SCALE_HASH_PW")
	case "ticket":
		hashPassword("TICKET_HASH_PW")
	case "":
		_, _ = fmt.Fprintln(os.Stderr, "error: service name not set (set Srvc variable in code)")
		os.Exit(1)
	default:
		if Srvc != "scale" && Srvc != "ticket" {
			_, _ = fmt.Fprintf(os.Stderr, "error: invalid service name '%s' (must be 'scale' or 'ticket')\n", Srvc)
			os.Exit(1)
		}
	}
}

func hashPassword(env string) {
	pw := strings.TrimSpace(os.Getenv(env))
	if pw == "" {
		// Empty password = no auth. Print empty string for ldflags.
		fmt.Print("")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// Base64 encode to avoid $ characters breaking ldflags
	fmt.Print(base64.StdEncoding.EncodeToString(hash))
}
