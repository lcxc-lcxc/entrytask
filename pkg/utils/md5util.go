package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// Hash 使用md5+盐值加密
func Hash(str string) string {
	hash := md5.Sum([]byte("@@83CX^&#)(" + str))
	//数组转切片 hash[:]
	return hex.EncodeToString(hash[:])
}

func HashVerify(hash, pwd string) bool {
	return Hash(pwd) == hash
}
