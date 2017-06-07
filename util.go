// Copyright 2017
// Author: catlittlechen@gmail.com

package nexus

import (
	"crypto/md5"
	"encoding/binary"
)

// MD5 .
func MD5(key string) [16]byte {
	return md5.Sum([]byte(key))
}

func hashValue(sum [16]byte, diff int) uint32 {
	if diff > 12 {
		return 0
	}
	return binary.BigEndian.Uint32(sum[diff : diff+4])
}
