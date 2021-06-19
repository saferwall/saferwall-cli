package sfwcli

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

const (
	// DefaultAuthUsername represents the key for reading the username
	// from env variables.
	DefaultAuthUsername = "SAFERWALL_AUTH_USERNAME"
	// DefaultAuthPassword represents the key for reading password
	// from env variables.
	DefaultAuthPassword = "SAFERWALL_AUTH_PASSWORD"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// SHA256 calculates the SHA256 hash of a given byte slice.
func SHA256(b []byte) string {
	h := sha256.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

// ExitWithError will trigger an print a message to stdout and exit.
func ExitWithError(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

// StringInSlice checks if a string is in a list.
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// SetAuthUsername writes username as an env variable.
func SetAuthUsername(key, value string) error {
	return os.Setenv(key, value)
}

// SetAuthPassword writes password as an env variable.
func SetAuthPassword(key, value string) error {
	return os.Setenv(key, value)
}

// SetAuthentificationData creates new env variables for authenticating users.
func SetAuthentificationData(username, password string) error {
	err := SetAuthUsername(DefaultAuthUsername, username)
	if err != nil {
		return err
	}
	return SetAuthPassword(DefaultAuthPassword, password)
}
