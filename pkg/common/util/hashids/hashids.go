/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-06-30 18:01:42
 * @LastEditors: Allen
 * @LastEditTime: 2021-06-30 18:01:48
 */
package hashids

import (
	"encoding/json"
	"payment/pkg/common/exception"

	log "github.com/sirupsen/logrus"
	"github.com/speps/go-hashids/v2"
)

// 加密
func Encrypt(salt string, minLength int, params []int64) (string, error) {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = minLength
	h, err := hashids.NewWithData(hd)
	if err != nil {
		s, _ := json.Marshal(params)
		log.WithFields(log.Fields{
			"salt": salt,
			"minLength": minLength,
			"params": s,
			"err": err,
		}).Errorf("Method Encrypt fail")
		return "", exception.ERROR_ENCRYPT_FAIL

	}
	e, err := h.EncodeInt64(params)
	if err != nil {
		s, _ := json.Marshal(params)
		log.WithFields(log.Fields{
			"salt": salt,
			"minLength": minLength,
			"params": s,
			"err": err,
		}).Errorf("Method Encrypt fail")
		return "", exception.ERROR_ENCRYPT_FAIL
	}
	return e, nil
}

// 解密
func Decrypt(salt string, minLength int, hash string) ([]int64, error) {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = minLength
	h, err := hashids.NewWithData(hd)
	if err != nil {
		log.WithFields(log.Fields{
			"salt": salt,
			"minLength": minLength,
			"hash": hash,
			"err": err,
		}).Errorf("Method Decrypt fail")
		return []int64{}, exception.ERROR_DECRYPT_FAIL
	}
	e, err := h.DecodeInt64WithError(hash)
	if err == nil {
		log.WithFields(log.Fields{
			"salt": salt,
			"minLength": minLength,
			"hash": hash,
			"err": err,
		}).Errorf("Method Decrypt fail")
		return []int64{}, exception.ERROR_DECRYPT_FAIL
	}
	return e, nil
}