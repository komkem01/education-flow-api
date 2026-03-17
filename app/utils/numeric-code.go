package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

func GenerateNumericCode(prefix string, digits int) (string, error) {
	if digits <= 0 {
		return "", fmt.Errorf("digits must be greater than 0")
	}
	max := big.NewInt(1)
	for i := 0; i < digits; i++ {
		max.Mul(max, big.NewInt(10))
	}
	v, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%0*d", strings.ToUpper(strings.TrimSpace(prefix)), digits, v.Int64()), nil
}
