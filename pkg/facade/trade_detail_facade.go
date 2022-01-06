/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-06-04 14:47:54
 * @LastEditors: Allen
 * @LastEditTime: 2021-06-04 14:47:55
 */
package facade

import (
	"payment/pkg/facade/dto"
)

/**
* @Description: facade interface
* @Author: Allen
 */
type TradeDetailFacade interface {
	Pay(req *dto.PayRequestDto) (string, error)

	Ack(req *dto.AckRequestDto) error

	NotifyAlipay(req *dto.AlipayNotifyDto) error

	NotifyWechat(req *dto.WechatNotifyInfoDto) error

}

