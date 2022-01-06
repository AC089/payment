/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-05-26 20:00:51
 * @LastEditors: Allen
 * @LastEditTime: 2021-05-26 20:04:52
 */
package repository

import (
	"gorm.io/gorm"
)

var (
	receiptReposiroty =new(receiptDetailReposiroty)
)

type receiptDetailReposiroty struct {
	db *gorm.DB
}

func ReceiptDetailReposirotyInstance() *receiptDetailReposiroty {
	return receiptReposiroty
}

/**
* @Description: set DB
* @Author: Allen
* @param {*gorm.DB} db
* @return {*}
* @error:
 */
func (p *receiptDetailReposiroty) SetDB(db *gorm.DB) *receiptDetailReposiroty {
	p.db = db
	return p
}


// func (p *receiptDetailReposiroty) Inject() *receiptDetailReposiroty {
	
// }