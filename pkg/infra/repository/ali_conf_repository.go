/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-05-26 11:46:19
 * @LastEditors: Allen
 * @LastEditTime: 2021-05-26 11:46:20
 */
package repository

import (
	"payment/pkg/domain/core"

	"gorm.io/gorm"
)

var (
	aliReposiroty = new(aliConfReposiroty)
)

type aliConfReposiroty struct {
	db *gorm.DB
}

func AliConfReposirotyInstance() *aliConfReposiroty {
	return aliReposiroty
}

/**
 * @Description: set DB
 * @Author: Allen
 * @param {*gorm.DB} db
 * @return {*}
 * @error: 
 */
func (a *aliConfReposiroty) SetDB(db *gorm.DB) *aliConfReposiroty {
	a.db = db
	return a
}

func (a *aliConfReposiroty) Inject() *aliConfReposiroty {
	core.AliConfInstance().SetDependency(a)
	return a
}



func (a *aliConfReposiroty)SelectListByCondition(condition core.AliConf) []core.AliConf {
	var aliConfs []core.AliConf
	a.db.Where(&condition).Find(&aliConfs)
	return aliConfs
}

func (a *aliConfReposiroty)SelectByCondition(condition core.AliConf) *core.AliConf {
	var aliConf *core.AliConf
	a.db.Where(&condition).First(aliConf)
	return aliConf
}