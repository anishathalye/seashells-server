package main

import (
	"math/rand"
	"strings"
)

const letters = "abcdefghjkmnpqrstuvwxyzABCDEFGHJKMNPQRSTUVWXYZ23456789"

func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func rSplitSingle(s, sep string) string {
	index := strings.LastIndex(s, sep)
	if index == -1 {
		return s
	}
	return s[:index]
}
