package utils

import (
	"strconv"
	"strings"
)

func ConvertRedisKeyToUintId(key string) (uint, error) {
	Uint64Id, err := strconv.ParseUint(key[strings.Index(key, ":")+1:len(key)], 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(Uint64Id), nil
}

func ConvertUintIdToRedisKey(prefix string, id uint) string {
	return prefix + ":" + strconv.FormatUint(uint64(id), 10)
}
