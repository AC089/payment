/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-06-05 16:24:31
 * @LastEditors: Allen
 * @LastEditTime: 2021-06-05 16:24:31
 */
package enums

import "payment/pkg/common/exception"

/**
 * @Description: 支付渠道
 * @Author: Allen
 */
type PayChannelEnum int32

const (
	ALIPAY_APP PayChannelEnum = 1 //支付宝APP
	ALIPAY_H5  PayChannelEnum = 2 //支付宝H5
	WECHAT_APP PayChannelEnum = 3 //微信支付APP
	WECHAT_H5  PayChannelEnum = 4 //微信支付H5
	IOSPAID    PayChannelEnum = 5 //IOS内购
)

func (obj PayChannelEnum) Code() int32 {
	return int32(obj)
}

func (obj PayChannelEnum) Desc() string {
	switch obj {
	case ALIPAY_APP:
		return "支付宝APP"
	case ALIPAY_H5:
		return "支付宝H5"
	case WECHAT_APP:
		return "微信支付APP"
	case WECHAT_H5:
		return "微信支付H5"
	case IOSPAID:
		return "IOS内购"
	default:
		return "UNKNOWN"
	}
}

/**
 * @Description: 配置文件选择策略
 * @Author: Allen
 */
type ConfStrategyEnum int32

const (
	POLLING ConfStrategyEnum = 1 //轮询
	RANDOM  ConfStrategyEnum = 2 //随机
	WEIGHT ConfStrategyEnum = 3 //权重
)

func (obj ConfStrategyEnum) Code() int32 {
	return int32(obj)
}

func (obj ConfStrategyEnum) Desc() string {
	switch obj {
	case POLLING:
		return "轮询"
	case RANDOM:
		return "随机"
	case WEIGHT:
		return "权重"
	default:
		return "UNKNOWN"
	}
}

/**
 * @Description: 支付宝密钥类型
 * @Author: Allen
 */
type SecretTypeEnum int32 

const (
	ORDINARY SecretTypeEnum = 1 //普通公私钥
	CERTIFICATE SecretTypeEnum = 2 //公私钥证书
)

func (obj SecretTypeEnum) Code() int32 {
	return int32(obj)
}

func (obj SecretTypeEnum) Desc() string {
	switch obj {
	case ORDINARY:
		return "普通公私钥"
	case CERTIFICATE:
		return "公私钥证书"
	default:
		return "UNKNOWN"
	}
}

/**
 * @Description: 签名类型
 * @Author: Allen
 */
type SignTypeEnum int32 

const (
	RSA SignTypeEnum = 1 //普通公私钥
	RSA2 SignTypeEnum = 2 //公私钥证书
)

func (obj SignTypeEnum) Code() int32 {
	return int32(obj)
}

func (obj SignTypeEnum) Desc() string {
	switch obj {
	case RSA:
		return "RSA"
	case RSA2:
		return "RSA2"
	default:
		return "UNKNOWN"
	}
}

func GetEnumByDesc(desc string) (SignTypeEnum, error) {
	switch desc {
	case "RSA":
		return RSA, nil
	case "RSA2":
		return RSA2, nil
	default:
		return 0, exception.ERROR_UNKNOWN
	}
}

/**
 * @Description: 数据状态枚举
 * @Author: Allen
 */
type StateEnum int32 

const (
	DISABLE StateEnum = 0
	ENABLE StateEnum = 1
)

func (obj StateEnum) Code() int32 {
	return int32(obj)
}

func (obj StateEnum) Desc() string {
	switch obj {
	case DISABLE:
		return "停用"
	case ENABLE:
		return "启用"
	default:
		return "UNKNOWN"
	}
}

/**
 * @Description: 交易状态枚举
 * @Author: Allen
 */
type StatusEnum int32 

const (
	UNPAID StatusEnum = 0
	PAID StatusEnum = 1
	BUSINESS_SUCCESS StatusEnum = 2
	CLOSE StatusEnum = 3
)

func (obj StatusEnum) Code() int32 {
	return int32(obj)
}

func (obj StatusEnum) Desc() string {
	switch obj {
	case UNPAID:
		return "未支付"
	case PAID:
		return "已支付"
	case BUSINESS_SUCCESS:
		return "业务处理成功"
	case CLOSE:
		return "交易关闭"
	default:
		return "UNKNOWN"
	}
}

/**
 * @Description: 支付宝交易状态枚举
 * @Author: Allen
*/
type AlipayTradeStatusEnum string 

const (
	WAIT_BUYER_PAY AlipayTradeStatusEnum = "WAIT_BUYER_PAY"
	TRADE_CLOSED AlipayTradeStatusEnum = "TRADE_CLOSED"
	TRADE_SUCCESS AlipayTradeStatusEnum = "TRADE_SUCCESS"
	TRADE_FINISHED AlipayTradeStatusEnum = "TRADE_FINISHED"
)

func (obj AlipayTradeStatusEnum) Code() string {
	return string(obj)
}

func (obj AlipayTradeStatusEnum) Desc() string {
	switch obj {
	case WAIT_BUYER_PAY:
		return "交易创建，等待买家付款。"
	case TRADE_CLOSED:
		return "未付款交易超时关闭，或支付完成后全额退款。"
	case TRADE_SUCCESS:
		return "交易支付成功。"
	case TRADE_FINISHED:
		return "交易结束，不可退款。"
	default:
		return "UNKNOWN"
	}
}


//微信回调状态
//SUCCESS：支付成功
// REFUND：转入退款
// NOTPAY：未支付
// CLOSED：已关闭
// REVOKED：已撤销（付款码支付）
// USERPAYING：用户支付中（付款码支付）
// PAYERROR：支付失败(其他原因，如银行返回失败)
type WechatTradeStatusEnum string 

const (
	WECHAT_SUCCESS WechatTradeStatusEnum = "SUCCESS"
	WECHAT_REFUND WechatTradeStatusEnum = "REFUND"
	WECHAT_NOTPAY WechatTradeStatusEnum = "NOTPAY"
	WECHAT_CLOSED WechatTradeStatusEnum = "CLOSED"
	WECHAT_REVOKED WechatTradeStatusEnum = "REVOKED"
	WECHAT_USERPAYING WechatTradeStatusEnum = "USERPAYING"
	WECHAT_PAYERROR WechatTradeStatusEnum = "PAYERROR"
)

func (obj WechatTradeStatusEnum) Code() string {
	return string(obj)
}

func (obj WechatTradeStatusEnum) Desc() string {
	switch obj {
	case WECHAT_SUCCESS:
		return "支付成功"
	case WECHAT_REFUND:
		return "转入退款"
	case WECHAT_NOTPAY:
		return "未支付"
	case WECHAT_CLOSED:
		return "已关闭"
	case WECHAT_REVOKED:
		return "已撤销（付款码支付）"
	case WECHAT_USERPAYING:
		return "用户支付中（付款码支付）"
	case WECHAT_PAYERROR:
		return "支付失败(其他原因，如银行返回失败)"
	default:
		return "UNKNOWN"
	}
}


/**
 * @Description: AckCode枚举
 * @Author: Allen
 */
type AckCodeEnum int32 

const (
	FAIL AckCodeEnum = 0
	SUCCESS AckCodeEnum = 1
)

func (obj AckCodeEnum) Code() int32 {
	return int32(obj)
}

func (obj AckCodeEnum) Desc() string {
	switch obj {
	case FAIL:
		return "业务执行错误，需要执行回滚退款"
	case SUCCESS:
		return "业务执行成功"
	default:
		return "UNKNOWN"
	}
}