package utils

import (
	"crypto/sha256"
)

var tronBase58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

func IsValidTRC20Address(address string) bool {
	if len(address) != 34 {
		return false
	}
	if address[0] != 'T' {
		return false
	}

	decoded, ok := base58Decode(address)
	if !ok {
		return false
	}
	if len(decoded) != 25 {
		return false
	}
	if decoded[0] != 0x41 {
		return false
	}

	payload := decoded[:21]
	checksum := decoded[21:]
	check := doubleSHA256(payload)[:4]
	return checksum[0] == check[0] && checksum[1] == check[1] && checksum[2] == check[2] && checksum[3] == check[3]
}

func doubleSHA256(b []byte) []byte {
	h1 := sha256.Sum256(b)
	h2 := sha256.Sum256(h1[:])
	out := make([]byte, 32)
	copy(out, h2[:])
	return out
}

func base58Decode(input string) ([]byte, bool) {
	if input == "" {
		return nil, false
	}

	index := make(map[byte]int)
	for i, c := range tronBase58Alphabet {
		index[c] = i
	}

	zeros := 0
	for zeros < len(input) && input[zeros] == '1' {
		zeros++
	}

	num := []byte{0}
	for i := 0; i < len(input); i++ {
		v, ok := index[input[i]]
		if !ok {
			return nil, false
		}
		num = mulAddBase58(num, v)
	}

	out := make([]byte, zeros+len(num))
	copy(out[zeros:], num)
	return out, true
}

func mulAddBase58(num []byte, add int) []byte {
	carry := add
	for i := len(num) - 1; i >= 0; i-- {
		val := int(num[i])*58 + carry
		num[i] = byte(val & 0xff)
		carry = val >> 8
	}
	for carry > 0 {
		num = append([]byte{byte(carry & 0xff)}, num...)
		carry >>= 8
	}

	i := 0
	for i < len(num)-1 && num[i] == 0 {
		i++
	}
	return num[i:]
}
