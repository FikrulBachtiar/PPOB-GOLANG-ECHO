package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func GenerateMD5(value string) string {
	hashMD5 := md5.New();
	hashMD5.Write([]byte(value));
	hashString := hashMD5.Sum(nil);

	return hex.EncodeToString(hashString);
}