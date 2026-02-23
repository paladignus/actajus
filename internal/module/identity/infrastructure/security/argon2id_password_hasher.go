// Package security
package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

type Argon2idPasswordHasher struct {
	Memory      uint32
	Time        uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

func NewArgon2idPasswordHasher() *Argon2idPasswordHasher {
	return &Argon2idPasswordHasher{
		Memory:      64 * 1024,
		Time:        3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
}

func (h *Argon2idPasswordHasher) Hash(plain string) (string, error) {
	if plain == "" {
		return "", errors.New("password is empty")
	}
	salt := make([]byte, h.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	key := argon2.IDKey([]byte(plain), salt, h.Time, h.Memory, h.Parallelism, h.KeyLength)
	saltB64 := base64.RawStdEncoding.EncodeToString(salt)
	keyB64 := base64.RawStdEncoding.EncodeToString(key)
	encoded := fmt.Sprintf(
		"argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, h.Memory, h.Time, h.Parallelism, saltB64, keyB64,
	)
	return encoded, nil
}

func (h *Argon2idPasswordHasher) Compare(encodedHash string, plain string) error {
	if encodedHash == "" || plain == "" {
		return errors.New("invalid credentials")
	}
	params, salt, expectedKey, err := decodeArgon2idHash(encodedHash)
	if err != nil {
		return errors.New("invalid credentials")
	}
	key := argon2.IDKey([]byte(plain), salt, params.Time, params.Memory, params.Parallelism, uint32(len(expectedKey)))
	if subtle.ConstantTimeCompare(key, expectedKey) != 1 {
		return errors.New("invalid credentials")
	}
	return nil
}

type argon2Params struct {
	Memory      uint32
	Time        uint32
	Parallelism uint8
}

func decodeArgon2idHash(encoded string) (argon2Params, []byte, []byte, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 5 {
		return argon2Params{}, nil, nil, errors.New("invalid hash format")
	}
	if parts[0] != "argon2id" {
		return argon2Params{}, nil, nil, errors.New("invalid algorithm")
	}
	if !strings.HasPrefix(parts[1], "v=") {
		return argon2Params{}, nil, nil, errors.New("invalid version")
	}
	verStr := strings.TrimPrefix(parts[1], "v=")
	ver, err := strconv.Atoi(verStr)
	if err != nil || ver != argon2.Version {
		return argon2Params{}, nil, nil, errors.New("unsupported version")
	}
	paramPart := parts[2]
	paramKV := strings.Split(paramPart, ",")
	if len(paramKV) != 3 {
		return argon2Params{}, nil, nil, errors.New("invalid params")
	}
	var p argon2Params
	for _, kv := range paramKV {
		kvp := strings.SplitN(kv, "=", 2)
		if len(kvp) != 2 {
			return argon2Params{}, nil, nil, errors.New("invalid params")
		}
		k := kvp[0]
		v := kvp[1]
		switch k {
		case "m":
			n, err := strconv.ParseUint(v, 10, 32)
			if err != nil {
				return argon2Params{}, nil, nil, errors.New("invalid memory")
			}
			p.Memory = uint32(n)
		case "t":
			n, err := strconv.ParseUint(v, 10, 32)
			if err != nil {
				return argon2Params{}, nil, nil, errors.New("invalid time")
			}
			p.Time = uint32(n)
		case "p":
			n, err := strconv.ParseUint(v, 10, 8)
			if err != nil {
				return argon2Params{}, nil, nil, errors.New("invalid parallelism")
			}
			p.Parallelism = uint8(n)
		default:
			return argon2Params{}, nil, nil, errors.New("invalid params")
		}
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil || len(salt) < 8 {
		return argon2Params{}, nil, nil, errors.New("invalid salt")
	}
	key, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil || len(key) < 16 {
		return argon2Params{}, nil, nil, errors.New("invalid key")
	}
	return p, salt, key, nil
}

func (h *Argon2idPasswordHasher) NeedsRehash(encodedHash string) bool {
	params, _, key, err := decodeArgon2idHash(encodedHash)
	if err != nil {
		return true
	}
	if params.Memory != h.Memory || params.Time != h.Time || params.Parallelism != h.Parallelism {
		return true
	}
	if uint32(len(key)) != h.KeyLength {
		return true
	}
	return false
}
