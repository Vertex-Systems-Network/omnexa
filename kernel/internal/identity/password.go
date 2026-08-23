package identity

import (
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
)

const (
	hashAlgorithm           = "pbkdf2-sha256"
	passwordIterations      = 600000
	passwordSaltBytes       = 16
	passwordDerivedKeyBytes = 32
	maxPasswordBytes        = 1024
	maxEncodedPasswordBytes = 512
)

// PasswordHasher is the narrow adaptive one-way password boundary consumed by
// P02.04. Implementations must never use reversible password storage.
type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(password, encoded string) (bool, error)
}

// PBKDF2PasswordHasher uses the Go standard-library PBKDF2-HMAC-SHA256 primitive
// with the governed 600,000-iteration work factor and a unique random salt.
type PBKDF2PasswordHasher struct{}

// NewPBKDF2PasswordHasher returns the governed password hasher.
func NewPBKDF2PasswordHasher() PBKDF2PasswordHasher { return PBKDF2PasswordHasher{} }

// Hash returns a PHC-style, self-describing one-way password representation.
func (PBKDF2PasswordHasher) Hash(password string) (string, error) {
	if !validPassword(password) {
		return "", passwordInvalidFailure()
	}
	salt := make([]byte, passwordSaltBytes)
	if _, err := rand.Read(salt); err != nil {
		return "", passwordHashFailure(err)
	}
	derived, err := pbkdf2.Key(sha256.New, password, salt, passwordIterations, passwordDerivedKeyBytes)
	if err != nil {
		return "", passwordHashFailure(err)
	}
	return fmt.Sprintf(
		"$%s$i=%d$%s$%s",
		hashAlgorithm,
		passwordIterations,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(derived),
	), nil
}

// Verify parses a bounded stored representation and compares the derived value
// in constant time. It never returns the password or stored hash in an error.
func (PBKDF2PasswordHasher) Verify(password, encoded string) (bool, error) {
	if !validPassword(password) || len(encoded) == 0 || len(encoded) > maxEncodedPasswordBytes {
		return false, nil
	}
	iterations, salt, expected, ok := parsePasswordHash(encoded)
	if !ok {
		return false, passwordHashInvalidFailure()
	}
	derived, err := pbkdf2.Key(sha256.New, password, salt, iterations, len(expected))
	if err != nil {
		return false, passwordHashFailure(err)
	}
	return secureEqual(derived, expected), nil
}

func parsePasswordHash(encoded string) (int, []byte, []byte, bool) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 5 || parts[0] != "" || parts[1] != hashAlgorithm || !strings.HasPrefix(parts[2], "i=") {
		return 0, nil, nil, false
	}
	iterations, err := strconv.Atoi(strings.TrimPrefix(parts[2], "i="))
	if err != nil || iterations < passwordIterations || iterations > passwordIterations*4 {
		return 0, nil, nil, false
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil || len(salt) < passwordSaltBytes || len(salt) > 64 {
		return 0, nil, nil, false
	}
	expected, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil || len(expected) != passwordDerivedKeyBytes {
		return 0, nil, nil, false
	}
	return iterations, salt, expected, true
}

func validPassword(password string) bool {
	return len(password) > 0 && len(password) <= maxPasswordBytes
}
