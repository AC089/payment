/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-06-23 20:45:59
 * @LastEditors: Allen
 * @LastEditTime: 2021-06-23 20:46:02
 */
package valueobject

import "time"


type AlipayNotifyVo struct {
	NotifyTime time.Time `json:"notify_time" form:"notify_time" time_format:"2006-01-02 15:04:05"`
	NotifyType string `json:"notify_type" form:"notify_type"`
	NotifyId string `json:"notify_id" form:"notify_id"`
	AppId string `json:"app_id" form:"app_id"`
	Charset string `json:"charset" form:"charset"`
	Version string `json:"version" form:"version"`
	SignType string `json:"sign_type" form:"sign_type"`
	Sign string `json:"sign" form:"sign"`
	TradeNo string `json:"trade_no" form:"trade_no"`
	OutTradeNo string `json:"out_trade_no" form:"out_trade_no"`
	OutBizNo string `json:"out_biz_no" form:"out_biz_no"`
	BuyerId string `json:"buyer_id" form:"buyer_id"`
	BuyerLogonId string `json:"buyer_logon_id" form:"buyer_logon_id"`
	SellerId string `json:"seller_id" form:"seller_id"`
	SellerEmail string `json:"seller_email" form:"seller_email"`
	TradeStatus string `json:"trade_status" form:"trade_status"`
	TotalAmount float64 `json:"total_amount" form:"total_amount"`
	ReceiptAmount float64 `json:"receipt_amount" form:"receipt_amount"`
	InvoiceAmount float64 `json:"invoice_amount" form:"invoice_amount"`
	BuyerPayAmount float64 `json:"buyer_pay_amount" form:"buyer_pay_amount"`
	PointAmount float64 `json:"point_amount" form:"point_amount"`
	RefundFee float64 `json:"refund_fee" form:"refund_fee"`
	Subject string `json:"subject" form:"subject"`
	Body string `json:"body" form:"body"`
	GmtCreate time.Time `json:"gmt_create" form:"gmt_create" time_format:"2006-01-02 15:04:05"`
	GmtPayment time.Time `json:"gmt_payment" form:"gmt_payment" time_format:"2006-01-02 15:04:05"`
	GmtRefund time.Time `json:"gmt_refund" form:"gmt_refund" time_format:"2006-01-02 15:04:05"`
	GmtClose time.Time `json:"gmt_close" form:"gmt_close" time_format:"2006-01-02 15:04:05"`
	FundBillList string `json:"fund_bill_list" form:"fund_bill_list"`
	PassbackParams string `json:"passback_params" form:"passback_params"`
	VoucherDetailList string `json:"voucher_detail_list" form:"voucher_detail_list"`
}