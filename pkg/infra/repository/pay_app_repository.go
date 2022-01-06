/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-05-26 19:53:54
 * @LastEditors: Allen
 * @LastEditTime: 2021-05-26 19:53:55
 */
package repository

import (
	"payment/pkg/domain/core"

	"gorm.io/gorm"
)
 
var (
	appReposiroty = new(payAppReposiroty)
)
 
type payAppReposiroty struct {
	db *gorm.DB
}
 
func PayAppReposirotyInstance() *payAppReposiroty {
	return appReposiroty
}

/**
* @Description: set DB
* @Author: Allen
* @param {*gorm.DB} db
* @return {*}
* @error:
*/
func (p *payAppReposiroty) SetDB(db *gorm.DB) *payAppReposiroty {
	p.db = db
	return p
}

func (p *payAppReposiroty) Inject() *payAppReposiroty {
	core.PayAppInstance().SetDependency(p)
	return p
}

func (p *payAppReposiroty) SelectByCondition(condition core.PayApp) *core.PayApp {
	var payApp core.PayApp
	p.db.Where(&condition).First(&payApp)
	return &payApp
}