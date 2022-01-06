/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-05-27 10:15:28
 * @LastEditors: Allen
 * @LastEditTime: 2021-05-27 10:15:29
 */
package core

import (
	"math/rand"
	"payment/pkg/common/enums"
	"payment/pkg/common/exception"
	"time"

	"gorm.io/gorm"
)

var(
	aliConf = new(AliConf)
	confStrategy = enums.RANDOM
	record = make(map[uint64]int)
)

type AliConfDependency interface {

	SelectByCondition(condition AliConf) *AliConf

	SelectListByCondition(condition AliConf) []AliConf
}


/**
 * @Description: 支付宝配置
 * @Author: Allen
 */
type AliConf struct {
	gorm.Model

	PayAppId uint `json:"payAppId"`
	AppId string `json:"appId"` //支付宝分配给开发者的应用ID
	MchId uint64 `json:"mchId"` //支付宝账户ID
	AppName string `json:"appName"` //应用名称
	PayChannel uint32 `json:"payChannel"` //支付渠道(1:alipay-app 2:alipay-h5)
	ReturnUrl string `json:"returnUrl"` //前台回跳地址
	NotifyUrl string `json:"notifyUrl"` //服务端异步通知地址
	PrivateKey string `json:"privateKey"` //应用私钥
	AlipayPublicKey string `json:"alipayPublicKey"` //支付宝公钥
	SignType uint32 `json:"signType"` //商户生成签名字符串所使用的签名算法类型,1:RSA2, 2:RSA
	Gateway string `json:"gateway"` //应用网关
	SecretType uint32 `json:"secretType"` //密钥类型, 1:普通公私钥, 2:公私钥证书
	AppCertPath string `json:"appCertPath"` //应用公钥证书绝对路径
	AlipayCertPath string `json:"alipayCertPath"` //支付宝公钥证书文件绝对路径
	AlipayRootCertPath string `json:"alipayRootCertPath"` //支付宝 CA 根证书文件绝对路径

	aliConfDependency AliConfDependency
}


func AliConfInstance() *AliConf {
	return aliConf
}


func (a *AliConf) SetDependency(dep AliConfDependency) *AliConf {
	a.aliConfDependency = dep
	return a
}

func (AliConf) TableName() string {
	return "ali_conf"
}


func (a *AliConf) SearchByPayAppId(payAppId uint, payChannel enums.PayChannelEnum) (*AliConf, error) {
	confs := a.aliConfDependency.SelectListByCondition(
		AliConf{PayAppId: payAppId, PayChannel: uint32(payChannel)})
	if len(confs) == 0 {
		return nil, exception.ERROR_DB_NIL
	}
	return a.strategy(payAppId, confs), nil
}

/**
 * @Description: 配置获取策略
 * @Author: Allen
 * @param {[]AliConf} confs
 * @return {*}
 * @error: 
 */
func (a *AliConf) strategy(payAppId uint, confs []AliConf) *AliConf {
	size := len(confs)
	if size == 1 {
		return &confs[0]
	}
	switch confStrategy {
	
	case enums.RANDOM:
		rand.Seed(time.Now().UnixNano())
		return &confs[rand.Intn(size - 1)]
	default:
		return &confs[0]
	}
}

func (a *AliConf) SearchByAppIdAndPayAppId(appId string, payAppId uint, payChannel uint32) (*AliConf, error) {
	conf := a.aliConfDependency.SelectByCondition(
		AliConf{PayAppId: payAppId, AppId: appId, PayChannel: payChannel})
	if conf.ID == 0 {
		return nil, exception.ERROR_DB_NIL
	}
	return conf, nil
}