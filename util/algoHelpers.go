package util

import (
	"crypto/md5"
	"math/big"
)

var Base62Alphabet = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

func Encode(i int) string {
	base := len(Base62Alphabet)
	var digits []int
	for i > 0 {
		r := i % base
		digits = append([]int{r}, digits...)
		i = i / base
	}

	var r []rune
	for _, d := range digits {
		r = append(r, Base62Alphabet[d])
	}
	return string(r)
}

func Hash(s string) int {
	hash := md5.Sum([]byte(s))
	hashInt := new(big.Int).SetBytes(hash[:6])

	return int(hashInt.Int64())
}
