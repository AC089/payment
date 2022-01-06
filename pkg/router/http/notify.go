/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-06-21 14:54:21
 * @LastEditors: Allen
 * @LastEditTime: 2021-06-21 14:54:25
 */
package http

import (
	"payment/pkg/common/constant"
	"payment/pkg/common/exception"
	"payment/pkg/common/util"
	"payment/pkg/common/util/hashids"
	"payment/pkg/facade/dto"
	"payment/pkg/infra/sdk/wechat"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var wehchatClient = wechat.GetInstance()

func (h HttpServer) NotifyAlipay (c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.WithFields(log.Fields{
				"err": err,
				"stack": util.CallStack(),
			}).Panicf("Method execution panic")
		}
	}()
	req := new(dto.AlipayNotifyDto)
	err := c.ShouldBind(req)
	if err != nil {
		log.Errorf("RestApi /notify/alipay is error, err:%v", err)
		c.String(200, "fail")
		return
	}
	err = h.tradeDetailFacade.NotifyAlipay(req)
	if err != nil && err != exception.ERROR_DUPLICATE_NOTIFY {
		log.Errorf("RestApi /notify/alipay is error, err:%v", err)
		c.String(200, "fail")
		return
	}
	c.String(200, "success")
}


func (h HttpServer) NotifyWechat (c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.WithFields(log.Fields{
				"err": err,
				"stack": util.CallStack(),
			}).Panicf("Method execution panic")
		}
	}()
	response := make(map[string]string)
	response["code"] = "SUCCESS"
	response["message"] = "成功"
	confId := c.Param("confId")
	confIds, err := hashids.Decrypt(constant.TRADE_ID_SALT, constant.TRADE_ID_MIN_LENGTH, confId)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Errorf("Method NotifyWechat is error")
		response["code"] = "Fail"
		response["message"] = "参数解析失败"
	}
	wechatCache, ok := wehchatClient.Cache.Load(uint64(confIds[0]))
	if !ok {
		log.WithFields(log.Fields{
			"err": err,
		}).Errorf("Method NotifyWechat is error")
		response["code"] = "Fail"
		response["message"] = "参数解密失败"
	}
	cache := wechatCache.(*wechat.WechatCache)
	successNotifyDto := new(dto.WechatSuccessNotifyDto)
	nre, err := cache.Handler.ParseNotifyRequest(*cache.Ctx, c.Request, successNotifyDto)
	if err != nil {
		log.Errorf("RestApi /notify/wechat is error, err:%v", err)
		response["code"] = "Fail"
		response["message"] = "解密验签失败"
		c.JSON(500, response)
		return
	}
	req := new(dto.WechatNotifyInfoDto)
	err = util.Assemble(nre, req)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Errorf("Method NotifyWechat is error")
		response["code"] = "Fail"
		response["message"] = "参数解析失败"
	}
	req.Resource.NotifyInfo = *successNotifyDto
	err = h.tradeDetailFacade.NotifyWechat(req)
	if err != nil && err != exception.ERROR_DUPLICATE_NOTIFY {
		log.Errorf("RestApi /notify/wechat is error, err:%v", err)
		response["code"] = "Fail"
		response["message"] = "业务处理失败"
		c.JSON(500, response)
		return
	}
	c.JSON(200, response)
}

