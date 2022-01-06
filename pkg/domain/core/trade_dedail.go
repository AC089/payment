/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-05-27 10:17:08
 * @LastEditors: Allen
 * @LastEditTime: 2021-05-27 10:17:09
 */
package core

import (
	"encoding/json"
	"fmt"
	"payment/pkg/common/exception"
	"payment/pkg/common/util"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var(
	tradeDetail = new(TradeDetail)
)

type TradeDetailDependency interface {

	Insert(trade *TradeDetail) int64

	SelectByCondition(trade TradeDetail) *TradeDetail

	UpdateById(trade TradeDetail) int64
}

type TradeDetail struct{
	gorm.Model

	PayAppId uint `json:"payAppId"` //发起支付appid
	PayAppCode string `json:"payAppCode"` //app代号
	ServerId string `json:"serverId"` //发起支付服务id
	PayChannel uint32 `json:"payChannel"` //支付渠道(1:alipay-app 2:alipay-h5 3:wechat-app 4:wechat-h5 5:iosPaid)
	ChanAppId string `json:"chanAppId"` //支付渠道应用id
	OutTradeNo string `json:"outTradeNo"` //商户App唯一订单号
	TransactionNo string `json:"transactionNo"` //支付渠道系统生成唯一订单号
	AccountId uint64 `json:"accountId"` //支付账户ID
	Amount int32 `json:"amount"` //支付金额,单位分
	ActualAmount int32 `json:"actualAmount"` //实际支付金额,单位分
	Subject string `json:"subject"` //商品标题
	Body string `json:"body"` //商品描述信息
	ClientIp string `json:"clientIp"` //客户端ip
	Device string `json:"device"` //设备标识
	Player string `json:"player"` //用户标识:alipay为支付宝唯一用户号,wechat为openid
	CreateTime uint32 `json:"createTime"` //交易创建时间,取回调通知返回值
	SuccessTime uint32 `json:"successTime"` //交易成功时间,取回调通知时间
	Status uint32 `json:"status"` //支付状态:0-未支付 1-支付成功 2-业务处理成功

	tradeDetailDependency TradeDetailDependency
}


func TradeDetailInstance() *TradeDetail {
	return tradeDetail
}

func (TradeDetail) TableName() string {
	return "trade_detail"
}

func (a *TradeDetail) SetDependency(dep TradeDetailDependency) *TradeDetail {
	a.tradeDetailDependency = dep
	return a
}

func (a *TradeDetail)CreateTradeDetail(trade *TradeDetail) error {
	trade.ID = uint(util.WorkerInstance().Next())
	//生成订单号
	trade.OutTradeNo = fmt.Sprintf("%d%d", util.GetYMDHMSInt(time.Now().Unix()), trade.ID)
	row := a.tradeDetailDependency.Insert(trade)
	if row < 1 {
		return exception.ERROR_DB_FAIL
	}
	return nil
}

/**
 * @Description: 根据tradeId查询交易记录
 * @Author: Allen
 * @param {uint} tradeId
 * @return {*}
 * @error: 
 */
func (a *TradeDetail) SearchById(tradeId uint) (*TradeDetail, error) {
	tradeDetail := TradeDetail{}
	tradeDetail.ID = tradeId
	result := a.tradeDetailDependency.SelectByCondition(tradeDetail)
	if result.ID == 0 {
		log.WithFields(log.Fields{
			"tradeId": tradeId,
			}).Errorf("method SearchById fail")
		return nil, exception.ERROR_DB_NIL
	}
	return result, nil
}


/**
 * @Description: 根据订单查询交易记录
 * @Author: Allen
 * @return {*}
 * @error: 
 */
func (a *TradeDetail) SearchByOutTradeNo(outTradeNo string) (*TradeDetail, error){
	tradeDetail := TradeDetail{}
	tradeDetail.OutTradeNo = outTradeNo
	result := a.tradeDetailDependency.SelectByCondition(tradeDetail)
	if result.ID == 0 {
		log.WithFields(log.Fields{
			"outTradeNo": outTradeNo,
			}).Errorf("method SearchByOutTradeNo fail")
		return nil, exception.ERROR_DB_NIL
	}
	return result, nil
}

/**
 * @Description: 更新交易状态
 * @Author: Allen
 * @return {*}
 * @error: 
 */
func (a *TradeDetail) UpdateTradeDetail(trade TradeDetail) error {
	row := a.tradeDetailDependency.UpdateById(trade)
	if row < 1 {
		jsonstr, _ := json.Marshal(trade)
		log.WithFields(log.Fields{
			"ID": trade.ID,
			"req": jsonstr, 
			}).Errorf("method UpdateTradeDetail fail")
		return exception.ERROR_DB_FAIL
	}
	return nil
}

