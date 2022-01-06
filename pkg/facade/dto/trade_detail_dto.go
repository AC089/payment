/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-06-05 16:00:21
 * @LastEditors: Allen
 * @LastEditTime: 2021-06-05 16:00:22
 */
package dto

import (
	"payment/pkg/common/enums"
)

/**
 * @Description: 发起支付dto
 * @Author: Allen
 * @param {*}
 * @return {*}
 * @error: 
 */
type PayRequestDto struct {
	PayAppCode string `json:"PayAppCode" validate:"required"`//app唯一代号
	ServerId string `json:"ServerId" validate:"required"`//发起支付服务id
	AccountId uint64 `json:"AccountId" validate:"required"`//订单号
	TimeExpire int64 `json:"TimeExpire" validate:"required"`//订单失效时间
	PayChannel enums.PayChannelEnum `json:"PayChannel" validate:"required"`//支付渠道(1:alipay-app 2:alipay-h5 3:wechat-app 4:wechat-h5 5:iosPaid)
	Amount uint64 `json:"Amount" validate:"required"`//支付金额
	Subject string `json:"Subject" validate:"required"`//商品标题
	Body string `json:"Body" validate:"required"`//商品标题
	ClientIp string `json:"ClientIp" validate:"required"`//客户端ip
	Device string `json:"Device" validate:"required"`//设备标识
}

/**
 * @Description: 通知客户端发放道具ACK结构体
 * @Author: Allen
 * @param {*}
 * @return {*}
 * @error: 
 */
type AckRequestDto struct {
	PayAppCode string   `json:"payAppCode"`
	OutTradeNo string   `json:"OutTradeNo"`
	AckCode enums.AckCodeEnum  `json:"AckCode"`
}