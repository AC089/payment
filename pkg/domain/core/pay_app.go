/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-05-27 10:16:34
 * @LastEditors: Allen
 * @LastEditTime: 2021-05-27 10:16:35
 */
package core

import (
	"payment/pkg/common/exception"

	"gorm.io/gorm"
)

var(
	payApp = new(PayApp)
)


//基础设施层仓储接口
type PayAppDependency interface {
	SelectByCondition(condition PayApp) *PayApp
	
}

/**
 * @Description: 应用配置
 * @Author: Allen
 */
type PayApp struct {
	gorm.Model
	
	PayAppCode string `json:"payAppCode"` //app唯一代号
	PayAppName string `json:"payAppName"` //app名称
	Description string `json:"description"` //描述

	payAppDependency PayAppDependency
}

func PayAppInstance() *PayApp {
	return payApp
}

func (a *PayApp) SetDependency(dep PayAppDependency) *PayApp {
	a.payAppDependency = dep
	return a
}

func (PayApp) TableName() string {
	return "pay_app"
}

/**
 * @Description: 根据app代码查询payApp
 * @Author: Allen
 * @param {string} code
 * @return {*}
 * @error: 
 */
 func (a *PayApp) SearchByCode(code string) (*PayApp, error) {
	result := a.payAppDependency.SelectByCondition(PayApp{PayAppCode: code})
	if result.ID == 0 {
		return nil, exception.ERROR_DB_NIL
	}
	return result, nil
}