package enc

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type Hasher struct {
	secret []byte
}

func NewHasher(secret string) *Hasher {
	return &Hasher{
		secret: []byte(secret),
	}
}

func (h *Hasher) Hash(token string) string {
	hh := hmac.New(sha256.New, h.secret)
	hh.Write([]byte(token))
	mac := hh.Sum(nil)

	return hex.EncodeToString(mac)
}

func (h *Hasher) Compare(hash string, val string) (bool, error) {
	h1 := hmac.New(sha256.New, h.secret)
	h1.Write([]byte(val))

	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		return false, err
	}

	return hmac.Equal(h1.Sum(nil), hashBytes), nil
}
