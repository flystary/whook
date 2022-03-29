package hash

import (
	"crypto/md5"
	"encoding/hex"
)


func Md5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}