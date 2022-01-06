/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-05-26 19:18:57
 * @LastEditors: Allen
 * @LastEditTime: 2021-05-26 19:53:57
 */
package repository

import (
	"payment/pkg/domain/core"

	"gorm.io/gorm"
)

var (
	wechatReposiroty = new(wechatConfReposiroty)
)

type wechatConfReposiroty struct {
	db *gorm.DB
}

func WechatConfReposirotyInstance() *wechatConfReposiroty {
	return wechatReposiroty
}

/**
 * @Description: set DB
 * @Author: Allen
 * @param {*gorm.DB} db
 * @return {*}
 * @error:
 */
func (w *wechatConfReposiroty) SetDB(db *gorm.DB) *wechatConfReposiroty {
	w.db = db
	return w
}



func (w *wechatConfReposiroty) Inject() *wechatConfReposiroty {
	core.WechatConfInstance().SetDependency(w)
	return w
}

func (w *wechatConfReposiroty) SelectByCondition(condition core.WechatConf) *core.WechatConf {
	var wechatConf *core.WechatConf
	w.db.Where(&condition).First(wechatConf)
	return wechatConf
}

func (w *wechatConfReposiroty) SelectListByCondition(condition core.WechatConf) []core.WechatConf {
	var wechatConfs []core.WechatConf
	w.db.Where(&condition).Find(&wechatConfs)
	return wechatConfs
}