package main

import (
	crand "crypto/rand"
	"fmt"
	"io"
	"math/rand"
	"time"
)

const _letters string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// NewUUID2 generates and returns a uuid
// @returns string
// @returns error
func NewUUID2() (string, error) {
	b := make([]byte, 16)
	n, err := io.ReadFull(crand.Reader, b)
	if n != len(b) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	b[8] = b[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	b[6] = b[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}

func Merge(l, r []int64) []int64 {
	ret := make([]int64, 0, len(l)+len(r))
	for len(l) > 0 || len(r) > 0 {
		if len(l) == 0 {
			return append(ret, r...)
		}
		if len(r) == 0 {
			return append(ret, l...)
		}
		if l[0] <= r[0] {
			ret = append(ret, l[0])
			l = l[1:]
		} else {
			ret = append(ret, r[0])
			r = r[1:]
		}
	}
	return ret
}

func MergeSort(s []int64) []int64 {
	if len(s) <= 1 {
		return s
	}
	n := len(s) / 2
	l := MergeSort(s[:n])
	r := MergeSort(s[n:])
	return Merge(l, r)
}
