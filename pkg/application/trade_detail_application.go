/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-06-04 14:31:47
 * @LastEditors: Allen
 * @LastEditTime: 2021-06-04 14:31:48
 */
package application

import (
	"encoding/json"
	"fmt"
	"payment/pkg/common/constant"
	"payment/pkg/common/enums"
	"payment/pkg/common/exception"
	"payment/pkg/common/util"
	"payment/pkg/common/util/hashids"
	"payment/pkg/domain/core"
	"payment/pkg/domain/valueobject"
	"payment/pkg/infra"
	"payment/pkg/infra/sdk/alipay"
	"payment/pkg/infra/sdk/alipay/kernel"
	"payment/pkg/infra/sdk/wechat"
	"time"

	log "github.com/sirupsen/logrus"
	rcore "github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/app"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
)

const (
	ALIPAY_GATEWAY = "https://openapi.alipay.com/gateway.do"
	CONTENT_TYPE = "application/json"
)

func NewTradeDetailApplication() *TradeDetailApplication {
	return &TradeDetailApplication{
		tradeDetail: core.TradeDetailInstance(),
		payApp: core.PayAppInstance(),
		aliConf: core.AliConfInstance(),
		wechatConf: core.WechatConfInstance(),
		alipayClient: alipay.GetInstance(),
		wechatClient: wechat.GetInstance(),
	}
}

type TradeDetailApplication struct {
	tradeDetail *core.TradeDetail
	payApp *core.PayApp
	aliConf *core.AliConf
	wechatConf *core.WechatConf
	alipayClient *alipay.AlipayClient
	wechatClient *wechat.WechatClient
}

/**
 * @Description: 支付宝app/h5支付
 * @Author: Allen
 * @param {*valueobject.PayRequestVo} payReq
 * @return {*}
 * @error: 
 */
func (p TradeDetailApplication) PayAliApp(payReq *valueobject.PayRequestVo) (string, error) {
	app, err := p.payApp.SearchByCode(payReq.PayAppCode)
	if err != nil {
		log.WithFields(log.Fields{"code": payReq.PayAppCode,"err": err,}).Errorf("App code incorrect")
		return "", err
	}
	conf, err := p.aliConf.SearchByPayAppId(app.ID, enums.PayChannelEnum(payReq.PayChannel))
	if err != nil {
		log.WithFields(log.Fields{"PayAppId": app.ID,"channel": enums.ALIPAY_APP, "err": err}).Errorf("aliConf incorrect")
		return "", err
	}
	//写入数据库
	tradeDetail := &core.TradeDetail{
		PayAppId: app.ID,
		PayAppCode: app.PayAppCode,
		ServerId: payReq.ServerId,
		PayChannel: conf.PayChannel,
		ChanAppId: conf.AppId,
		AccountId: payReq.AccountId,
		Amount: payReq.Amount,
		Subject: payReq.Subject,
		Body: payReq.Body,
		ClientIp: payReq.ClientIp,
		Device: payReq.Device,
		CreateTime: uint32(time.Now().Unix()),
		Status: uint32(enums.UNPAID),
	}
	err = p.tradeDetail.CreateTradeDetail(tradeDetail)
	if err != nil {
		log.Errorf("method PayAliApp error, err:%v", err)
		return "", err
	}
	client := p.alipayClient.GetClient(*conf)
	biz := valueobject.BizContent{
		TotalAmount: fmt.Sprintf("%.2f", float64(payReq.Amount)/float64(100)),
		Subject: payReq.Subject,
		OutTradeNo: tradeDetail.OutTradeNo,
		TimeExpire: util.TimeFormat(payReq.TimeExpire),
		NotifyUrl: conf.NotifyUrl,
	}
	m, err := util.StructToMapStr(biz)
	if err != nil {
		log.Errorf("StructToMapStr error, err:%v", err)
		return "", exception.ERROR_JSONPARSE
	}
	// //混淆订单号,拼接回调地址
	// hash, err := hashids.Encrypt(constant.TRADE_ID_SALT, constant.TRADE_ID_MIN_LENGTH, []int64{int64(tradeDetail.ID)})
	// if err != nil {
	// 	log.Errorf("method PayAliApp error, err:%v", err)
	// 	return "", err
	// }
	// m["notify_url"] = fmt.Sprintf("%s/%s", conf.NotifyUrl, hash)
	result, err := client.Execute("alipay.trade.app.pay", nil, m)
	if !alipay.IsSuccess(*result.Code) {
		log.WithFields(log.Fields{
			"code": result.Code,
			"msg": result.Msg,
			"subCode": result.SubCode,
			"subMsg": result.SubMsg,
			}).Errorf("Alipay alipay.trade.app.pay Execute Fail")
		return "", exception.ERROR_UNKNOWN
	}
	return *result.HttpBody, nil
}

/**
 * @Description: 微信App发起支付
 * @Author: Allen
 * @param {*valueobject.PayRequestVo} payReq
 * @return {*}
 * @error: 
 */
func (p TradeDetailApplication) PayWechatApp(payReq *valueobject.PayRequestVo) (string, error) {
	napp, err := p.payApp.SearchByCode(payReq.PayAppCode)
	if err != nil {
		log.WithFields(log.Fields{"code": payReq.PayAppCode,"err": err,}).Errorf("App code incorrect")
		return "", err
	}
	conf, err := p.wechatConf.SearchByPayAppId(napp.ID, enums.PayChannelEnum(payReq.PayChannel))
	if err != nil {
		log.WithFields(log.Fields{"PayAppId": napp.ID,"channel": enums.ALIPAY_APP, "err": err}).Errorf("aliConf incorrect")
		return "", err
	}
	hash, err := hashids.Encrypt(constant.TRADE_ID_SALT, constant.TRADE_ID_MIN_LENGTH, []int64{int64(conf.ID)})
	if err != nil {
		log.WithFields(log.Fields{"confId": conf.ID,"err": err,}).Errorf("hashids Encrypt fail")
		return "", err
	}
	//写入数据库
	tradeDetail := &core.TradeDetail{
		PayAppId: napp.ID,
		PayAppCode: napp.PayAppCode,
		ServerId: payReq.ServerId,
		PayChannel: conf.PayChannel,
		ChanAppId: conf.AppId,
		AccountId: payReq.AccountId,
		Amount: payReq.Amount,
		Subject: payReq.Subject,
		Body: payReq.Body,
		ClientIp: payReq.ClientIp,
		Device: payReq.Device,
		CreateTime: uint32(time.Now().Unix()),
		Status: uint32(enums.UNPAID),
	}
	err = p.tradeDetail.CreateTradeDetail(tradeDetail)
	if err != nil {
		log.Errorf("method PayAliApp error, err:%v", err)
		return "", err
	}
	cache, err := p.wechatClient.GetCache(*conf)
	if err != nil {
		log.Errorf("GetCache error, err:%v", err)
		return "", exception.ERROR_UNKNOWN
	}
	svc := app.AppApiService{Client: cache.Client}
	resp, _, err := svc.Prepay(*cache.Ctx, 
		app.PrepayRequest{
			Appid:       rcore.String(conf.AppId),
			Mchid:       rcore.String(conf.MahId),
			Description: rcore.String(payReq.Subject),
			OutTradeNo:  rcore.String(tradeDetail.OutTradeNo),
			// Attach:      rcore.String("自定义数据说明"),
			TimeExpire: rcore.Time(util.GetTime(payReq.TimeExpire)),
			NotifyUrl:   rcore.String(fmt.Sprintf("%s/%s", conf.NotifyUrl, hash)),
			Amount: &app.Amount{
				Total: rcore.Int32(int32(payReq.Amount)),
			},
		})
	if err != nil {
		log.Errorf("wechat App Prepay error, confId:%d, err:%v", conf.ID, err)
		return "", exception.ERROR_UNKNOWN
	}
	result, err := json.Marshal(resp)
	if err != nil {
		log.Errorf("Json Marshal error, err:%v", err)
		return "", exception.ERROR_JSONPARSE
	}
	return string(result), nil
}

/**
 * @Description: 微信H5发起支付
 * @Author: Allen
 * @param {*valueobject.PayRequestVo} payReq
 * @return {*}
 * @error: 
 */
func (p TradeDetailApplication) PayWechatH5(payReq *valueobject.PayRequestVo) (string, error) {
	napp, err := p.payApp.SearchByCode(payReq.PayAppCode)
	if err != nil {
		log.WithFields(log.Fields{"code": payReq.PayAppCode,"err": err,}).Errorf("App code incorrect")
		return "", err
	}
	conf, err := p.wechatConf.SearchByPayAppId(napp.ID, enums.PayChannelEnum(payReq.PayChannel))
	if err != nil {
		log.WithFields(log.Fields{"PayAppId": napp.ID,"channel": enums.ALIPAY_APP, "err": err}).Errorf("aliConf incorrect")
		return "", err
	}
	hash, err := hashids.Encrypt(constant.TRADE_ID_SALT, constant.TRADE_ID_MIN_LENGTH, []int64{int64(conf.ID)})
	if err != nil {
		log.WithFields(log.Fields{"confId": conf.ID,"err": err,}).Errorf("hashids Encrypt fail")
		return "", err
	}
	//写入数据库
	tradeDetail := &core.TradeDetail{
		PayAppId: napp.ID,
		PayAppCode: napp.PayAppCode,
		ServerId: payReq.ServerId,
		PayChannel: conf.PayChannel,
		ChanAppId: conf.AppId,
		AccountId: payReq.AccountId,
		Amount: payReq.Amount,
		Subject: payReq.Subject,
		Body: payReq.Body,
		ClientIp: payReq.ClientIp,
		Device: payReq.Device,
		CreateTime: uint32(time.Now().Unix()),
		Status: uint32(enums.UNPAID),
	}
	err = p.tradeDetail.CreateTradeDetail(tradeDetail)
	if err != nil {
		log.Errorf("method PayAliApp error, err:%v", err)
		return "", err
	}
	cache, err := p.wechatClient.GetCache(*conf)
	if err != nil {
		log.Errorf("GetCache error, err:%v", err)
		return "", exception.ERROR_UNKNOWN
	}
	svc := h5.H5ApiService{Client: cache.Client}
	resp, _, err := svc.Prepay(*cache.Ctx, 
		h5.PrepayRequest{
			Appid:       rcore.String(conf.AppId),
			Mchid:       rcore.String(conf.MahId),
			Description: rcore.String(payReq.Subject),
			OutTradeNo:  rcore.String(tradeDetail.OutTradeNo),
			// Attach:      rcore.String("自定义数据说明"),
			TimeExpire: rcore.Time(util.GetTime(payReq.TimeExpire)),
			NotifyUrl:   rcore.String(fmt.Sprintf("%s/%s", conf.NotifyUrl, hash)),
			Amount: &h5.Amount{
				Total: rcore.Int32(int32(payReq.Amount)),
			},
		})
	if err != nil {
		log.Errorf("wechat H5 Prepay error, confId:%d, err:%v", conf.ID, err)
		return "", exception.ERROR_UNKNOWN
	}
	result, err := json.Marshal(resp)
	if err != nil {
		log.Errorf("Json Marshal error, err:%v", err)
		return "", exception.ERROR_JSONPARSE
	}
	return string(result), nil
}

/**
 * @Description: 支付宝支付回调通知
 * @Author: Allen
 * @param {*valueobject.AlipayNotifyVo} req
 * @return {*}
 * @error: 
 */
func (p TradeDetailApplication) NotifyAlipay(req *valueobject.AlipayNotifyVo) error {
	//查询订单，确认通知有效
	trade, err := p.tradeDetail.SearchByOutTradeNo(req.OutTradeNo)
	if err != nil {
		log.WithFields(log.Fields{"outTradeNo": req.OutTradeNo, "err": err}).Errorf("method NotifyAlipay is error")
		return err
	}
	//校验订单
 	err = checkTrade(*trade, int32(req.TotalAmount * 100), req.AppId, req.OutTradeNo)
	if err != nil {
		log.WithFields(log.Fields{"OutTradeNo": req.OutTradeNo, "err": err}).Errorf("method NotifyAlipay is error")
		return err
	}
	//TODO 校验通知中的 seller_id（或者 seller_email ) 是否为 out_trade_no 这笔单据的对应的操作方（有的时候，一个商户可能有多个seller_id/seller_email），

	//查询支付宝配置获取平台公钥
	conf, err := p.aliConf.SearchByAppIdAndPayAppId(trade.ChanAppId, trade.PayAppId, trade.PayChannel)
	if err != nil {
		log.WithFields(log.Fields{"ChanAppId": trade.ChanAppId, "PayAppId": trade.PayAppId, "err": err}).
		Errorf("method NotifyAlipay is error")
		return err
	}
	//验证签名
	client := p.alipayClient.GetClient(*conf)
	params, err := util.StructToMap(req)
	if err != nil {
		log.Errorf("method NotifyAlipay error, err:%v", err)
		return err
	}
	delete(params, "sign")
	delete(params, "sign_type")
	content := util.MapTOStr(params)
	if !kernel.AsyncVerify(req.Sign, content, client.AlipayPublicKey) {
		log.Errorf("Sign fail, sign:%s, content:%s", req.Sign, content)
		return exception.ERROR_SIGN_FAIL
	}
	//验签成功修改数据状态
	td, err := fillAlipayTradeParams(trade.ID, *req)
	if err != nil {
		if err != nil {
			log.WithFields(log.Fields{"err": err}).
			Errorf("method NotifyAlipay is error")
			return err
		}
	}
	err = p.tradeDetail.UpdateTradeDetail(td)
	if err != nil {
		log.Errorf("UpdateTradeDetail fail, err:%v", err)
		return err
	}
	//校验订单状态，只有已支付才通知发放道具
	if td.Status == uint32(enums.PAID) {
		notifyDelivery(*trade)
	}
	return nil
}

/**
 * @Description: 微信支付回调
 * @Author: Allen
 * @param {*valueobject.WechatNotifyInfoVo} req
 * @return {*}
 * @error: 
 */
func (p TradeDetailApplication) NotifyWechat(req *valueobject.WechatNotifyInfoVo) error {
	//如果不是支付通知成功直接return
	if req.ResourceType != "TRANSACTION.SUCCESS" {
		return nil
	}
	//查询订单，确认通知有效
	trade, err := p.tradeDetail.SearchByOutTradeNo(req.Resource.NotifyInfo.OutTradeNo)
	if err != nil {
		log.WithFields(log.Fields{"outTradeNo": req.Resource.NotifyInfo.OutTradeNo, "err": err}).Errorf("method NotifyAlipay is error")
		return err
	}
	//校验订单
	err = checkTrade(*trade, int32(req.Resource.NotifyInfo.Amount.Total * 100), req.Resource.NotifyInfo.Appid, req.Resource.NotifyInfo.OutTradeNo)
	if err != nil {
		log.WithFields(log.Fields{"OutTradeNo": req.Resource.NotifyInfo.OutTradeNo, "err": err}).Errorf("method NotifyWechat is error")
		return err
	}
	//修改数据状态
	td, err := fillWechatReqParams(trade.ID, req.Resource.NotifyInfo)
	if err != nil {
		if err != nil {
			log.WithFields(log.Fields{"err": err}).
			Errorf("method NotifyWechat is error")
			return err
		}
	}
	err = p.tradeDetail.UpdateTradeDetail(td)
	if err != nil {
		log.Errorf("UpdateTradeDetail fail, err:%v", err)
		return err
	}
	//校验订单状态，只有已支付才通知发放道具
	if td.Status == uint32(enums.PAID) {
		notifyDelivery(*trade)
	}
	return nil
}

/**
 * @Description: 校验订单
 * @Author: Allen
 * @param {core.TradeDetail} trade
 * @param {int32} amount
 * @param {*} appId
 * @param {string} outTradeNo
 * @return {*}
 * @error: 
 */
func checkTrade(trade core.TradeDetail, amount int32, appId, outTradeNo string) error {
	//校验订单号是否无误
	if trade.OutTradeNo != outTradeNo {
		log.Errorf("订单号不一致,OutTradeNo:%s,tradeId:%d", trade.OutTradeNo, trade.ID)
		return exception.ERROR_INCONSISTENT_ORDER
	}
	//校验订单状态,过滤重复通知
	if trade.Status == uint32(enums.BUSINESS_SUCCESS) {
		log.Errorf("重复通知,OutTradeNo:%s", trade.OutTradeNo)
		return exception.ERROR_DUPLICATE_NOTIFY
	}
	//校验订单金额
	if trade.Amount != amount {
		log.Errorf("订单金额不匹配,amount:%d", amount)
		return exception.ERROR_INVALID_NOTIFY
	}
	//验证 app_id 是否为该商户本身
	if trade.ChanAppId != appId {
		log.Errorf("app_id不匹配,appId:%s", appId)
		return exception.ERROR_INVALID_NOTIFY
	}
	return nil
}


func fillWechatReqParams(tradeId uint, req valueobject.WechatSuccessNotifyVo) (core.TradeDetail, error) {
	trade := core.TradeDetail{}
	trade.ID = tradeId
	var status uint32
	switch enums.WechatTradeStatusEnum(req.TradeState) {
	case enums.WECHAT_NOTPAY:
		status = uint32(enums.UNPAID)
	case enums.WECHAT_CLOSED, enums.WECHAT_PAYERROR, enums.WECHAT_REFUND:
		status = uint32(enums.CLOSE)
	case enums.WECHAT_SUCCESS:
		status = uint32(enums.PAID)
	default:
		log.WithFields(log.Fields{
			"tradeStatus": req.TradeState, 
			}).Errorf("value tradeStatus fail")
		return trade, exception.ERROR_PARAMETER
	}
	trade.Status = status
	trade.TransactionNo = req.TransactionId
	trade.ActualAmount = int32(req.Amount.PayerTotal)
	nilTime := time.Time{}
	if nilTime != req.SuccessTime && status == uint32(enums.PAID) {
		notifyTime := req.SuccessTime.Unix()
		trade.SuccessTime = uint32(notifyTime)
	}
	trade.Player = req.Player.Openid
	return trade, nil
}


/**
 * @Description: 填充alilpay交易参数
 * @Author: Allen
 * @param {uint} tradeId
 * @param {*valueobject.AlipayNotifyVo} req
 * @return {*}
 * @error: 
 */
func fillAlipayTradeParams(tradeId uint, req valueobject.AlipayNotifyVo) (core.TradeDetail, error) {
	trade := core.TradeDetail{}
	trade.ID = tradeId
	var status uint32
	switch enums.AlipayTradeStatusEnum(req.TradeStatus) {
	case enums.WAIT_BUYER_PAY:
		status = uint32(enums.UNPAID)
	case enums.TRADE_CLOSED:
		status = uint32(enums.CLOSE)
	case enums.TRADE_SUCCESS, enums.TRADE_FINISHED:
		status = uint32(enums.PAID)
	default:
		log.WithFields(log.Fields{
			"tradeStatus": req.TradeStatus, 
			}).Errorf("value tradeStatus fail")
		return trade, exception.ERROR_PARAMETER
	}
	trade.Status = status
	trade.TransactionNo = req.TradeNo
	trade.ActualAmount = int32(req.ReceiptAmount * 100)
	nilTime := time.Time{}
	if nilTime != req.NotifyTime && status == uint32(enums.PAID) {
		notifyTime := req.NotifyTime.Unix()
		trade.SuccessTime = uint32(notifyTime)
	}
	trade.Player = req.BuyerId
	return trade, nil
}




/**
 * @Description: 通知发放道具
 * @Author: Allen
 * @param {TradeDetail} trade
 * @return {*}
 * @error: 
 */
func notifyDelivery(trade core.TradeDetail) error {
	//发送消息队列通知调用服务更新订单状态
	notify := valueobject.PaymentNotifyVo{
		NotifyTime: uint64(time.Now().Unix()),
		OutTradeNo: trade.OutTradeNo,
		SuccessTime: trade.SuccessTime,
		ReceiptAmount: trade.ActualAmount,
		BuyerId: trade.Player,
	}
	str, err := json.Marshal(notify)
	if err != nil {
		log.WithFields(log.Fields{
			"ID": trade.ID,
			"err": err, 
			}).Errorf("payment notifyDelivery json Marshal fail")
		return exception.ERROR_JSONPARSE
	}
	//发布消息
	topic := fmt.Sprintf("%s_%s_%s",infra.GetNsqInstance().TopicPre, trade.PayAppCode, trade.ServerId)
	infra.GetNsqInstance().Publish(topic, []byte(str))
	return nil
}

/**
 * @Description: 回调通知Ack
 * @Author: Allen
 * @param {*valueobject.AckRequestVo} req
 * @return {*}
 * @error: 
 */
func (p TradeDetailApplication) Ack(req *valueobject.AckRequestVo) error {
	//查询订单，确认订单有效
	trade, err := p.tradeDetail.SearchByOutTradeNo(req.OutTradeNo)
	if err != nil {
		log.WithFields(log.Fields{"OutTradeNo": req.OutTradeNo, "err": err}).Errorf("method Ack is error")
		return err
	}
	if trade.Status != uint32(enums.PAID) {
		log.WithFields(log.Fields{"OutTradeNo": req.OutTradeNo, "Status": trade.Status}).
		Errorf("订单状态必须为已支付")
		return exception.ERROR_TRADE_STATUS
	}
	//订单有效修改订单状态为业务处理成功
	td := core.TradeDetail{
		Status: uint32(enums.BUSINESS_SUCCESS),
	}
	td.ID = trade.ID
	err = p.tradeDetail.UpdateTradeDetail(td)
	if err != nil {
		log.Errorf("UpdateTradeDetail fail, err:%v", err)
		return err
	}
	return nil
}