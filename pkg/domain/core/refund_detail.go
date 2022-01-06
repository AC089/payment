/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-05-27 10:16:58
 * @LastEditors: Allen
 * @LastEditTime: 2021-05-27 10:16:59
 */
package core

import (
	"gorm.io/gorm"
)

var(
	refundDetail = new(RefundDetail)
)

type RefundDetailDependency interface {

}

/**
 * @Description: 退款详情
 * @Author: Allen
 */
type RefundDetail struct {
	gorm.Model

	TransactionNo string `json:"transactionNo"` //支付渠道系统生成唯一订单号
	OutTradeNo string `json:"outTradeNo"` //商户网站|App唯一订单号
	OutRefundNo string `json:"outRefundNo"` //退款订单号
	Reason string `json:"reason"` //退款原因
	TradeDetailId uint `json:"tradeDetailId"` //交易详情id
	Amount int32 `json:"amount" ` //退款金额,单位分
	Status uint32 `json:"status"` //退款状态,0-退款处理中 1-退款成功 2-退款关闭 3-退款异常
	ClientIp string `json:"clientIp"` //客户端ip
	Device string `json:"device"` //设备标识
	CreateTime uint32 `json:"createTime"` //退款创建时间,取回调通知返回值
	SuccessTime uint32 `json:"successTime"` //退款成功时间,取回调通知时间


	refundDetailDependency RefundDetailDependency
}

func (RefundDetail) TableName() string {
	return "refund_detail"
}

func RefundDetailInstance() *RefundDetail {
	return refundDetail
}

func (a *RefundDetail) SetDependency(dep RefundDetailDependency) *RefundDetail {
	a.refundDetailDependency = dep
	return a
}