/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-06-15 11:11:56
 * @LastEditors: Allen
 * @LastEditTime: 2021-06-15 11:12:07
 */
package util

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"payment/pkg/common/enums"
	"payment/pkg/common/exception"
	"sort"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)


/**
 * @Description: 签名公共函数
 * @Author: Allen
 * @param {map[string]interface{}} params
 * @param {string} key
 * @param {enums.SignTypeEnum} signType
 * @return {*}
 * @error: 
 */
func Sign(params map[string]interface{}, key string, signType enums.SignTypeEnum) (string, error) {
	signContent := MapTOStr(params)
	var hash crypto.Hash
	switch signType {
	case enums.RSA:
		hash = crypto.SHA1
	case enums.RSA2:
		hash = crypto.SHA256
	default:
		log.WithFields(log.Fields{
			"signType": signType,
		}).Error("枚举参数错误")
		return "", exception.ERROR_PARAMETER
	}
	return RsaSign(signContent, key, hash)
}

/**
 * @Description: 结构对象转map
 * @Author: Allen
 * @param {interface{}} s
 * @return {*}
 * @error: 
 */
 func StructToMapStr(s interface{}) (map[string]string, error) {
	params := make(map[string]string)
	str, err := json.Marshal(s)
	if err != nil {
		log.Errorf("JSON Marshal error, err:%v", err)
		return params, exception.ERROR_JSONPARSE
	}
	err = json.Unmarshal([]byte(str), &params)
	if err != nil {
		log.Errorf("JSON Marshal error, err:%v", err)
		return params, exception.ERROR_JSONPARSE
	}
	return params, nil
}

/**
 * @Description: 结构对象转map
 * @Author: Allen
 * @param {interface{}} s
 * @return {*}
 * @error: 
 */
func StructToMap(s interface{}) (map[string]interface{}, error) {
	params := make(map[string]interface{})
	str, err := json.Marshal(s)
	if err != nil {
		log.Errorf("JSON Marshal error, err:%v", err)
		return params, exception.ERROR_JSONPARSE
	}
	err = json.Unmarshal([]byte(str), &params)
	if err != nil {
		log.Errorf("JSON Marshal error, err:%v", err)
		return params, exception.ERROR_JSONPARSE
	}
	return params, nil
}

/**
 * @Description: map按key排序并转string
 * @Author: Allen
 * @param {map[string]string} params
 * @return {*}
 * @error: 
 */
func MapTOStr(params map[string]interface{}) string {
	//ksort
    var keys []string
    for k := range params {
        keys = append(keys, k)
    }
    sort.Strings(keys)
	 //拼接
	var dataParams string
	for _, k := range keys {
		if IsEmptyOrNil(params[k]) {
			continue
		}
		log.Debugf("key:%s  value:%s", k, params[k])
		dataParams = dataParams + k + "=" + Strval(params[k]) + "&"
    }
	//去掉最后一个&
	dataParams = dataParams[0 : len(dataParams)-1]
	return dataParams
}

// Strval 获取变量的字符串值
// 浮点型 3.0将会转换成字符串3, "3"
// 非数值或字符类型的变量将会被转换成JSON格式字符串
func Strval(value interface{}) string {
	// interface 转 string
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	case time.Time:
		v := value.(time.Time)
		key = TimeFormat(v.Unix())
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}


/**
 * @Description: Sha256WithRSA\ShaWithRSA 签名算法
 * @Author: Allen
 * @param {string} signContent
 * @param {string} privateKey
 * @param {crypto.Hash} hash
 * @return {*}
 * @error: 
 */
func RsaSign(signContent string, privateKey string, hash crypto.Hash) (string, error) {
	shaNew := hash.New()
	shaNew.Write([]byte(signContent))
	hashed := shaNew.Sum(nil)
	priKey, err := ParsePrivateKey(privateKey)
	if err != nil {
		return "", err
	}

	signature, err := rsa.SignPKCS1v15(rand.Reader, priKey, hash, hashed)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

func ParsePrivateKey(privateKey string)(*rsa.PrivateKey, error) {
	// 2、解码私钥字节，生成加密对象
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		return nil, exception.ERROR_PRIVATE_KEY
	}
	// 3、解析DER编码的私钥，生成私钥对象
	priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return priKey, nil
}
