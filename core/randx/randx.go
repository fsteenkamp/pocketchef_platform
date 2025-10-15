package randx

import (
	"crypto/rand"
	"math/big"
)

const PINLen = 6
const TokenLen = 32

func String(n int, options string) string {
	ret := make([]byte, n)
	for i := range n {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(options))))
		if err != nil {
			panic(err)
		}
		ret[i] = options[num.Int64()]
	}

	return string(ret)
}

const AlphaNumeric = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const Numbers = "0123456789"

func UID() string {
	return String(32, AlphaNumeric)
}

// Pin generates and returns a random 6 character Pin.
func PIN() string {
	return String(PINLen, Numbers)
}
