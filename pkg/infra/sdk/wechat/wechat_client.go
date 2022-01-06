/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-06-21 20:16:44
 * @LastEditors: Allen
 * @LastEditTime: 2021-06-21 20:16:44
 */
package wechat

import (
	"context"
	"payment/pkg/infra/httpclient"
	"sync"

	ncore "payment/pkg/domain/core"

	log "github.com/sirupsen/logrus"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
)

var (
	instance = new(WechatClient)
	customHTTPClient = httpclient.GetInstance()
)

type WechatCache struct {
	Ctx *context.Context //上下文
	Handler *notify.Handler //回调处理器
	Client *core.Client //sdk客户端
}

type WechatClient struct {
	Cache sync.Map //微信支付sdk客户端缓存
	mutex *sync.Mutex
}

func GetInstance() *WechatClient {
	return instance
}

func (p WechatClient) GetCache(conf ncore.WechatConf) (*WechatCache, error) {
	cache, ok := p.Cache.Load(conf.ID)
	if ok {
		return cache.(*WechatCache), nil
	}
	init, err := p.InitClient(conf)
	if err != nil {
		log.Errorf("初始化微信支付sdk失败")
		return nil, err
	}
	//存入缓存
	p.Cache.Store(conf.ID, init)
	return init, nil
}

func (p WechatClient) InitClient(conf ncore.WechatConf) (*WechatCache, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	//双重检查缓存
	cli, ok := p.Cache.Load(conf.ID)
	if ok {
		return cli.(*WechatCache), nil
	}
	log.Infof("初始化Wehchat配置,confId:%d", conf.ID)
	ctx := context.Background()
	//加载私钥
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(conf.PrivateKeyPath)
	if err != nil {
		log.Errorf("私钥证书加载失败")
		return nil, err
	}
	opts := []core.ClientOption{
		// 一次性设置 签名/验签/敏感字段加解密，并注册 平台证书下载器，自动定时获取最新的平台证书
		option.WithWechatPayAutoAuthCipher(conf.MahId, conf.Serial, mchPrivateKey, conf.Apiv3),
		// 设置自定义 HTTPClient 实例，不设置时默认使用 http.Client{}，并设置超时时间为 30s
		option.WithHTTPClient(customHTTPClient),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		log.Printf("new wechat pay client err:%s", err)
		return nil, err
	}
	// 由于已经使用 WithWechatPayAutoAuthCipher 在 downloader.MgrInstance() 中注册了商户的下载器，则下面可以直接使用该下载管理器获取证书
	certVisitor := downloader.MgrInstance().GetCertificateVisitor(conf.MahId)
	handler := notify.NewNotifyHandler(conf.Apiv3, verifiers.NewSHA256WithRSAVerifier(certVisitor))
	return &WechatCache{
		Ctx : &ctx,
		Handler: handler,
		Client: client,
	}, nil
}