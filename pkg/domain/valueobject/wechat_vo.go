/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-06-23 20:46:06
 * @LastEditors: Allen
 * @LastEditTime: 2021-06-23 20:46:06
 */
package valueobject

import "time"

/**
 * @Description: 微信支付回调Vo
 * @Author: Allen
 * @param {*}
 * @return {*}
 * @error:
 */
type WechatNotifyInfoVo struct {
	ID           string             `json:"id"`
	CreateTime   *time.Time         `json:"create_time"`
	EventType    string             `json:"event_type"`
	ResourceType string             `json:"resource_type"`
	Resource     *ResourceVo 		`json:"resource"`
	Summary      string             `json:"summary"`
}

type ResourceVo struct {
	Algorithm      string `json:"algorithm"`
	Ciphertext     string `json:"ciphertext"`
	AssociatedData string `json:"associated_data"`
	Nonce          string `json:"nonce"`
	OriginalType   string `json:"original_type"`

	Plaintext string // Ciphertext 解密后内容

	NotifyInfo WechatSuccessNotifyVo // 解析内容后结构体
}


type WechatSuccessNotifyVo struct {
	Appid string `json:"appid"`
	Mchid string `json:"mchid"`
	OutTradeNo string `json:"out_trade_no"`
	TransactionId string `json:"transaction_id"`
	TradeType string `json:"trade_type"`
	TradeState string `json:"trade_state"`
	SuccessTime time.Time `json:"success_time"`
	Player PlayerVo `json:"player"`
	Amount AmountVo `json:"amount"`
}

type PlayerVo struct {
	Openid string `json:"openid"`
}

type AmountVo struct {
	Total uint64 `json:"total"`
	PayerTotal uint64 `json:"payer_total"`
}