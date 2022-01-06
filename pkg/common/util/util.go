/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-05-28 15:36:06
 * @LastEditors: Allen
 * @LastEditTime: 2021-05-28 15:36:14
 */
package util

import (
	"bytes"
	"os"
	"os/exec"
	"payment/pkg/common/exception"
	"reflect"
	"runtime"
	"sync"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

/**
 * @Description: 执行shell命令
 * @Author: Allen
 * @param {string} s
 * @return {*}
 * @error:
 */
func Exec_shell(s string) (string, error) {
	cmd := exec.Command("/bin/bash", "-c", s)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	return out.String(), err
}


func Mkdir(dir string) bool {
	// 创建文件夹
	err := os.Mkdir(dir, os.ModePerm)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("mkdir failed!")
		return false
	}
	return true
}

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

/**
 * @Description: 数据装配
 * @Author: Allen
 * @param {interface{}} source
 * @param {interface{}} target
 * @return {*}
 * @error: 
 */
func Assemble(source interface{}, target interface{}) error {
	// src和dst必须是struct，且dst必须是point
	srcType := reflect.TypeOf(source)
	srcValue := reflect.ValueOf(source)
	dstValue := reflect.ValueOf(target)

	if srcType.Kind() != reflect.Struct {
		return exception.ERROR_TYPE_ASSEMBLE
	}
	if dstValue.Kind() != reflect.Ptr {
		return exception.ERROR_TYPE_ASSEMBLE
	}

	for i := 0; i < srcType.NumField(); i++ {
		dstField :=  dstValue.Elem().FieldByName(srcType.Field(i).Name)
		if dstField.CanSet() {
			dstField.Set(srcValue.Field(i))
		}
	}

	return nil
}

func IsEmptyOrNil(obj interface{}) bool {
	if obj == nil {
		return true
	}
	switch obj.(type) {
	case string:
		if obj.(string) == "" {
			return true
		}
	case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
		if obj == 0 {
			return true
		}
	default:
		return false
	}
	return false
}

func CallStack() string {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	return string(buf[:n])
}

/**
 * @Description: 返回参数校验实例
 * @Author: Allen
 * @param {*}
 * @return {*}
 * @error: 
 */
var doOnce sync.Once
func Validator() *validator.Validate {
	var instance *validator.Validate 
	doOnce.Do(func() {
		instance = validator.New()	
	})
	return instance
}