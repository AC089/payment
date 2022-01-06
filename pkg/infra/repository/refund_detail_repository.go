/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-05-26 20:04:50
 * @LastEditors: Allen
 * @LastEditTime: 2021-05-26 20:04:50
 */
package repository

import (
	"payment/pkg/domain/core"

	"gorm.io/gorm"
)

var (
	refundReposiroty = new(refundDetailReposiroty)
)

type refundDetailReposiroty struct {
	db *gorm.DB
}

func RefundDetailReposirotyInstance() *refundDetailReposiroty {
	return refundReposiroty
}

/**
* @Description: set DB
* @Author: Allen
* @param {*gorm.DB} db
* @return {*}
* @error:
*/
func (p *refundDetailReposiroty) SetDB(db *gorm.DB) *refundDetailReposiroty {
	p.db = db
	return p
}

func (p *refundDetailReposiroty) Inject() *refundDetailReposiroty {
	core.RefundDetailInstance().SetDependency(p)
	return p
}

func (p *refundDetailReposiroty) SelectList() {
	
}