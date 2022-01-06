/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-06-23 20:46:38
 * @LastEditors: Allen
 * @LastEditTime: 2021-06-23 20:46:40
 */
package dto

import "time"

/**
 * @Description: 微信支付回调结构体
 * @Author: Allen
 * @param {*}
 * @return {*}
 * @error:
 */
type WechatNotifyInfoDto struct {
	ID           string             `json:"id"`
	CreateTime   *time.Time         `json:"create_time"`
	EventType    string             `json:"event_type"`
	ResourceType string             `json:"resource_type"`
	Resource     *ResourceDto 		`json:"resource"`
	Summary      string             `json:"summary"`
}

type ResourceDto struct {
	Algorithm      string `json:"algorithm"`
	Ciphertext     string `json:"ciphertext"`
	AssociatedData string `json:"associated_data"`
	Nonce          string `json:"nonce"`
	OriginalType   string `json:"original_type"`

	Plaintext string // Ciphertext 解密后内容

	NotifyInfo WechatSuccessNotifyDto // 解析内容后结构体
}


type WechatSuccessNotifyDto struct {
	Appid string `json:"appid"`
	Mchid string `json:"mchid"`
	OutTradeNo string `json:"out_trade_no"`
	TransactionId string `json:"transaction_id"`
	TradeType string `json:"trade_type"`
	TradeState string `json:"trade_state"`
	SuccessTime *time.Time `json:"success_time"`
	Player PlayerDto `json:"player"`
	Amount AmountDto `json:"amount"`
}

type PlayerDto struct {
	Openid string `json:"openid"`
}

type AmountDto struct {
	Total uint64 `json:"total"`
	PayerTotal uint64 `json:"payer_total"`
}