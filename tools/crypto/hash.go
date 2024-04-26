package crypto

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"strings"

	"golang.org/x/crypto/argon2"
)

type Hash interface {
	Generate(string) (string, error)
	Compare(string, string) (bool, error)
}

type argon2ID struct {
	params *GeneratePwdParams
}

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

// GeneratePwdParams is generate password params data type
type GeneratePwdParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

// NewArgon2IDHash return argon2id hash
func NewArgon2IDHash(params *GeneratePwdParams) Hash {
	return &argon2ID{params: params}
}

// Generate generates hash from password
func (aid *argon2ID) Generate(password string) (encodedHash string, err error) {
	// Generate a cryptographically secure random salt.
	p := aid.params
	salt, err := generateRandomBytes(p.SaltLength)
	if err != nil {
		return "", err
	}

	hashed := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hashed)

	// Return a string using the standard encoded hash representation.
	encodedHash = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, p.Memory, p.Iterations, p.Parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Compare compares password with encoded hash
func (aid *argon2ID) Compare(password, encodedHash string) (match bool, err error) {
	// Extract the parameters, salt and derived key from the encoded password
	// hash.
	p, salt, hashed, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Derive the key from the other password using the same parameters.
	otherHash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(hashed, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func decodeHash(encodedHash string) (p *GeneratePwdParams, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	elementLen := 6
	if len(vals) != elementLen {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	p = &GeneratePwdParams{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Iterations, &p.Parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.SaltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.KeyLength = uint32(len(hash))

	return p, salt, hash, nil
}

// HMAC generate hash with specific algorithm and key
func HMAC(hashFunc func() hash.Hash, secret []byte, data string) string {
	h := hmac.New(hashFunc, secret)
	h.Write([]byte(data))
	hashed := hex.EncodeToString(h.Sum(nil))

	return hashed
}
