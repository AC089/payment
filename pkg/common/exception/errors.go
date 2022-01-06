/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-06-07 17:14:47
 * @LastEditors: Allen
 * @LastEditTime: 2021-06-07 17:14:47
 */
package exception

import (
	"fmt"
	"strconv"
	"strings"
)


//业务异常
type JsonParseError int
func (p JsonParseError) Error() string { return fmt.Sprintf("%d:%s", p, "json数据解析错误") }

type PrivateKeyError int
func (p PrivateKeyError) Error() string { return fmt.Sprintf("%d:%s", p, "私钥错误") }

type ParameterError int
func (p ParameterError) Error() string { return fmt.Sprintf("%d:%s", p, "参数错误") }

type TypeAssembleError int
func (p TypeAssembleError) Error() string { return fmt.Sprintf("%d:%s", p, "数据装配失败") }

type PayChannelError int
func (p PayChannelError) Error() string { return fmt.Sprintf("%d:%s", p, "支付渠道错误") }

type TradeStatusError int
func (p TradeStatusError) Error() string { return fmt.Sprintf("%d:%s", p, "订单状态错误") }

type DuplicateTradeError int
func (p DuplicateTradeError) Error() string { return fmt.Sprintf("%d:%s", p, "订单重复") }

type DuplicateNotifyError int
func (p DuplicateNotifyError) Error() string { return fmt.Sprintf("%d:%s", p, "通知重复") }

type SignFailError int
func (p SignFailError) Error() string { return fmt.Sprintf("%d:%s", p, "验签失败") }

type InvalidNotify int
func (p InvalidNotify) Error() string { return fmt.Sprintf("%d:%s", p, "无效通知") }

type DecryptNotify int
func (p DecryptNotify) Error() string { return fmt.Sprintf("%d:%s", p, "解密通知数据失败") }

type InconsistentOrder int
func (p InconsistentOrder) Error() string { return fmt.Sprintf("%d:%s", p, "订单号不一致") }

type EncryptFail int
func (p EncryptFail) Error() string { return fmt.Sprintf("%d:%s", p, "id混淆加密失败") }

type DecryptFail int
func (p DecryptFail) Error() string { return fmt.Sprintf("%d:%s", p, "id混淆解密失败") }

//db错误
type DbFAILError int
func (p DbFAILError) Error() string { return fmt.Sprintf("%d:%s", p, "数据库操作失败") }
type DbNilError int
func (p DbNilError) Error() string { return fmt.Sprintf("%d:%s", p, "查询数据库返回空") }
//未知错误
type UnknownError int
func (p UnknownError) Error() string { return fmt.Sprintf("%d:%s", p, "未知错误") }


const (
	ERROR_JSONPARSE = JsonParseError(100)
	ERROR_PARAMETER = ParameterError(101)
	ERROR_TYPE_ASSEMBLE = TypeAssembleError(102)
	ERROR_PRIVATE_KEY = PrivateKeyError(103)
	ERROR_PAY_CHANNEL = PayChannelError(104)
	ERROR_TRADE_STATUS = TradeStatusError(105)
	ERROR_DUPLICATE_TRADE = DuplicateTradeError(106)
	ERROR_DUPLICATE_NOTIFY = DuplicateNotifyError(107)
	ERROR_SIGN_FAIL = SignFailError(108)
	ERROR_INVALID_NOTIFY = InvalidNotify(109)
	ERROR_DECRYPT_NOTIFY = DecryptNotify(100)
	ERROR_INCONSISTENT_ORDER = InconsistentOrder(101)
	ERROR_ENCRYPT_FAIL = EncryptFail(102)
	ERROR_DECRYPT_FAIL = DecryptFail(103)

	ERROR_DB_FAIL = DbFAILError(601)
	ERROR_DB_NIL = DbNilError(602)


	ERROR_UNKNOWN = UnknownError(999)
	
)


func ErrCode(err error) int32 {
	s := strings.Split(err.Error(), ":")
	if len(s) == 0 {
		return int32(ERROR_UNKNOWN)
	}
	i, err := strconv.Atoi(s[0])
	if err != nil {
		return int32(ERROR_UNKNOWN)
	}
	return int32(i)
}

func ErrMsg(err error) string {
	s := strings.Split(err.Error(), ":")
	if len(s) < 2 {
		return err.Error()
	}
	return s[1]
}
