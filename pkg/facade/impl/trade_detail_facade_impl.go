/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-06-05 17:06:56
 * @LastEditors: Allen
 * @LastEditTime: 2021-06-05 17:06:57
 */
package impl

import (
	"payment/pkg/application"
	"payment/pkg/common/enums"
	"payment/pkg/common/exception"
	"payment/pkg/common/util"
	"payment/pkg/domain/valueobject"
	"payment/pkg/facade/dto"

	log "github.com/sirupsen/logrus"
)


func NewTradeDetailFacadeImpl() *TradeDetailFacadeImpl {
	return &TradeDetailFacadeImpl{tradeDetailApplication: application.NewTradeDetailApplication()}
}

type TradeDetailFacadeImpl struct {
	tradeDetailApplication *application.TradeDetailApplication
}

/**
 * @Description: 支付接口
 * @Author: Allen
 * @param {*dto.PayRequestDto} req
 * @return {*}
 * @error: 
 */
func (p *TradeDetailFacadeImpl) Pay(req *dto.PayRequestDto) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			log.WithFields(log.Fields{
				"err": err,
				"stack": util.CallStack(),
			}).Panicf("Method execution panic")
		}
	}()
	//参数校验
	err := util.Validator().Struct(req)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Errorf("Method Pay validate param is error")
		return "", exception.ERROR_PARAMETER
	}
	payRequest := new(valueobject.PayRequestVo)
	err = util.Assemble(req, payRequest)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Errorf("Method Pay is error")
		return "", exception.ERROR_TYPE_ASSEMBLE
	}
	var result string
	switch req.PayChannel {
	case enums.ALIPAY_APP, enums.ALIPAY_H5:
		result, err = p.tradeDetailApplication.PayAliApp(payRequest)
	case enums.WECHAT_APP:
		result, err = p.tradeDetailApplication.PayWechatApp(payRequest)
	case enums.WECHAT_H5:
		result, err = p.tradeDetailApplication.PayWechatH5(payRequest)
	default:
		log.Errorf("PayChannel is error, PayChannel:%v", req.PayChannel)
		return "", exception.ERROR_PAY_CHANNEL
	}
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Errorf("Method Pay is error")
		return "", err
	}
	return result, nil
}

/**
 * @Description: 回调成功ack接口
 * @Author: Allen
 * @param {*dto.AckRequestDto} req
 * @return {*}
 * @error: 
 */
func (p *TradeDetailFacadeImpl) Ack(req *dto.AckRequestDto) error {
	defer func() {
		if err := recover(); err != nil {
			log.WithFields(log.Fields{
				"err": err,
				"stack": util.CallStack(),
			}).Panicf("Method execution panic")
		}
	}()
	//参数校验
	err := util.Validator().Struct(req)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Errorf("Method Ack validate param is error")
		return exception.ERROR_PARAMETER
	}
	ackRequest := new(valueobject.AckRequestVo)
	err = util.Assemble(req, ackRequest)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Errorf("Method Ack is error")
		return exception.ERROR_TYPE_ASSEMBLE
	}
	err = p.tradeDetailApplication.Ack(ackRequest)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Errorf("Method Ack is error")
		return err
	}
	return nil
}


/**
 * @Description: alipay回调通知
 * @Author: Allen
 * @param {*dto.AlipayNotifyDto} req
 * @return {*}
 * @error: 
 */
func (p*TradeDetailFacadeImpl)NotifyAlipay(req *dto.AlipayNotifyDto) error {
	
	
	vo := new(valueobject.AlipayNotifyVo)
	err := util.Assemble(req, vo)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Errorf("Method NotifyAlipay is error")
		return exception.ERROR_TYPE_ASSEMBLE
	}
	err = p.tradeDetailApplication.NotifyAlipay(vo)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Errorf("Method NotifyAlipay is error")
		return err
	}
	return nil
}


/**
 * @Description: wechat回调通知
 * @Author: Allen
 * @param {*dto.WechatNotifyInfoDto} req
 * @return {*}
 * @error: 
 */
func (p*TradeDetailFacadeImpl)NotifyWechat(req *dto.WechatNotifyInfoDto) error {
	defer func() {
		if err := recover(); err != nil {
			log.WithFields(log.Fields{
				"err": err,
				"stack": util.CallStack(),
			}).Panicf("Method execution panic")
		}
	}()
	
	vo := new(valueobject.WechatNotifyInfoVo)
	err := util.Assemble(req, vo)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Errorf("Method NotifyWechat is error")
		return exception.ERROR_TYPE_ASSEMBLE
	}
	err = p.tradeDetailApplication.NotifyWechat(vo)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Errorf("Method NotifyWechat is error")
		return err
	}
	return nil
}