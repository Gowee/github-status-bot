package utils

import (
	"os"
)

func IsPathNotExisting(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

// Max returns the larger of x or y.
func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// Min returns the smaller of x or y.
func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

// Ref: https://yourbasic.org/golang/max-min-function/
