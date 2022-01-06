/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-05-24 19:19:29
 * @LastEditors: Allen
 * @LastEditTime: 2021-05-24 19:19:29
 */
package util

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}