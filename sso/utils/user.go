/**
 * @Author: Robby
 * @File name: user.go
 * @Create date: 2021-11-03
 * @Function:
 **/

package utils

import (
	"crypto/md5"
	"encoding/hex"
)

const secret = "Md5SortKey"

func EncryptPassword(password string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(password)))
}
