/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-05-27 10:17:18
 * @LastEditors: Allen
 * @LastEditTime: 2021-05-27 10:17:19
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
	wechatConf = new(WechatConf)
)

type WechatConfDependency interface {

	SelectByCondition(condition WechatConf) *WechatConf

	SelectListByCondition(condition WechatConf) []WechatConf
}

/**
 * @Description: 微信支付配置
 * @Author: Allen
 */
type WechatConf struct {
	gorm.Model

	PayAppId uint `json:"payAppId"`
	MahId string `json:"mahId"` //微信商户号ID
	AppId string `json:"appId"` //应用ID
	PayChannel uint32 `json:"payChannel"` //支付渠道(3:wechat-app 4:wechat-h5)
	ApiKey string `json:"apiKey"` //API密钥
	NotifyUrl string `json:"notifyUrl"` //服务端异步通知地址
	PrivateKeyPath string `json:"privateKeyPath"` //私钥证书绝对路径
	PlatCertPath string `json:"platCertPath"` //支付平台公钥证书绝对路径
	Serial string `json:"serial"`
	Apiv3 string `json:"apiv3"`

	wechatConfDependency WechatConfDependency
}

func WechatConfInstance() *WechatConf {
	return wechatConf
}

func (WechatConf) TableName() string {
	return "wechat_conf"
}

func (w *WechatConf) SetDependency(dep WechatConfDependency) *WechatConf {
	w.wechatConfDependency = dep
	return w
}

func (w *WechatConf) SearchByPayAppId(payAppId uint, payChannel enums.PayChannelEnum) (*WechatConf, error) {
	confs := w.wechatConfDependency.SelectListByCondition(
		WechatConf{PayAppId: payAppId, PayChannel: uint32(payChannel)})
	if len(confs) == 0 {
		return nil, exception.ERROR_DB_NIL
	}
	return w.strategy(payAppId, confs), nil
}

/**
 * @Description: 配置获取策略
 * @Author: Allen
 * @param {[]WechatConf} confs
 * @return {*}
 * @error: 
 */
func (w *WechatConf) strategy(payAppId uint, confs []WechatConf) *WechatConf {
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

func (w *WechatConf) SearchByAppIdAndPayAppId(appId string, payAppId uint, payChannel uint32) (*WechatConf, error) {
	conf := w.wechatConfDependency.SelectByCondition(
		WechatConf{PayAppId: payAppId, AppId: appId, PayChannel: payChannel})
	if conf.ID == 0 {
		return nil, exception.ERROR_DB_NIL
	}
	return conf, nil
}