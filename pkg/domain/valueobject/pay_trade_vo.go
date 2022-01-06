/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-06-08 16:50:01
 * @LastEditors: Allen
 * @LastEditTime: 2021-06-08 16:50:02
 */
package valueobject

import (
	"payment/pkg/common/enums"
)

/**
 * @Description: 支付宝业务数据结构体
 * @Author: Allen
 * @param {*}
 * @return {*}
 * @error:
 */
type BizContent struct {
	TotalAmount string `json:"total_amount"`
	Subject string `json:"subject"`
	OutTradeNo string `json:"out_trade_no"`
	TimeExpire string `json:"time_expire"`
	NotifyUrl string `json:"notify_url"`
}

type PayRequestVo struct {
	PayAppCode string `json:"PayAppCode"`//app唯一代号
	ServerId string `json:"ServerId"`//发起支付服务id
	AccountId uint64 `json:"AccountId"`//账户id
	TimeExpire int64 `json:"TimeExpire"`//订单失效时间
	PayChannel enums.PayChannelEnum `json:"PayChannel"`//支付渠道(1:alipay-app 2:alipay-h5 3:wechat-app 4:wechat-h5 5:iosPaid)
	Amount int32 `json:"Amount"`//支付金额
	Currency string `json:"Currency"`//cny：人民币，境内商户号仅支持人民币。
	Subject string `json:"Subject"`//商品标题
	Body string `json:"Body"`//商品标题
	ClientIp string `json:"ClientIp"`//客户端ip
	Device string `json:"Device"`//设备标识
}



type PaymentNotifyVo struct {
	NotifyTime uint64 `json:"notifyTime"` 
	AccountId uint64 `json:"AccountId"` //payment服务交易ID
	OutTradeNo string `json:"outTradeNo"`
	SuccessTime uint32 `json:"successTime"`
	ReceiptAmount int32 `json:"receiptAmount"` //用户实际付款金额
	BuyerId string `json:"buyerId"` //买家支付宝账号对应的支付宝唯一用户号/微信openid
}


type AckRequestVo struct {
	PayAppCode string   `json:"payAppCode"`
	OutTradeNo string   `json:"outTradeNo"`
	AckCode enums.AckCodeEnum  `json:"ackCode"`
}
