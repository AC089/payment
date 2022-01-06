/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-06-22 15:58:56
 * @LastEditors: Allen
 * @LastEditTime: 2021-06-22 15:58:57
 */
package alipay

import (
	"sync"

	"payment/pkg/common/enums"
	ncore "payment/pkg/domain/core"
	"payment/pkg/infra/sdk/alipay/kernel"

	log "github.com/sirupsen/logrus"
)

var (
	instance = new(AlipayClient)
)

type AlipayClient struct {

	ClientMap sync.Map //微信支付sdk客户端缓存
	mutex *sync.Mutex
}

func GetInstance() *AlipayClient {
	return instance
}

func (p AlipayClient) GetClient(conf ncore.AliConf) (*kernel.Client) {
	client, ok := p.ClientMap.Load(conf.ID)
	if ok {
		return client.(*kernel.Client)
	}
	init := p.InitClient(conf)
	//存入缓存
	p.ClientMap.Store(conf.ID, init)
	return init
}



func (p AlipayClient) InitClient(conf ncore.AliConf) (*kernel.Client) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	//双重检查缓存
	cli, ok := p.ClientMap.Load(conf.ID)
	if ok {
		return cli.(*kernel.Client)
	}
	log.Infof("初始化Alipay配置,confId:%d", conf.ID)
	account:= &kernel.Client{}
	account.Protocol = "https"
	account.GatewayHost = "openapi.alipay.com"
	account.AppId = conf.AppId
	account.SignType = enums.SignTypeEnum(conf.SignType).Desc()
	if conf.SecretType == uint32(enums.ORDINARY) {
		account.AlipayPublicKey = conf.AlipayPublicKey
		account.MerchantPrivateKey = conf.PrivateKey
	} else {
		account.MerchantPrivateKey = conf.PrivateKey
		account.AlipayCertPath = conf.AlipayCertPath
		account.AlipayRootCertPath = conf.AlipayRootCertPath
		account.MerchantCertPath = conf.AppCertPath
	}
	
	return account
}

func IsSuccess(code string) bool {
	if code != "10000" {
		return false
	}
	return true
}