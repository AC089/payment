/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-06-16 14:44:46
 * @LastEditors: Allen
 * @LastEditTime: 2021-06-16 14:44:49
 */
package util

import (
	"strconv"
	"time"
)

/**
 * @Description: 获取当前格式化时间
 * @Author: Allen
 * @param {*}
 * @return {*}
 * @error:
 */
func NowFormat() string {
	return TimeFormat(time.Now().Unix())
}

// 获得的format time
func TimeFormat(ts int64) string {
	now := time.Unix(ts, 0)
	local1, _ := time.LoadLocation("Asia/Shanghai")
	return now.In(local1).Format("2006-01-02 15:04:05")
}

//获取年月日时分秒20200202
func GetYMDHMSInt(ts int64) int64 {
	dateStr := time.Unix(ts, 0).Format("20060102150405")
	date, _ := strconv.ParseInt(dateStr, 10, 64)
	return date
}

//时间戳转time
func GetTime(ts int64) time.Time {
	tm := time.Unix(ts, 0)
	timeStr := tm.Format("2006-01-02 15:04:05")
	local1, _ := time.LoadLocation("Asia/Shanghai")
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr, local1)
	return t
}