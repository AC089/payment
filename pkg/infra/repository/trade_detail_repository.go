/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-05-26 19:56:43
 * @LastEditors: Allen
 * @LastEditTime: 2021-05-26 19:56:46
 */
package repository

import (
	"payment/pkg/domain/core"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)


var (
	tradeReposiroty = new(tradeDetailReposiroty)
)


type tradeDetailReposiroty struct {
	db *gorm.DB
}

func TradeDetailReposirotyInstance() *tradeDetailReposiroty {
	return tradeReposiroty
}

/**
* @Description: set DB
* @Author: Allen
* @param {*gorm.DB} db
* @return {*}
* @error:
*/
func (p *tradeDetailReposiroty) SetDB(db *gorm.DB) *tradeDetailReposiroty {
	p.db = db
	return p
}

func (p *tradeDetailReposiroty) Inject() *tradeDetailReposiroty {
	core.TradeDetailInstance().SetDependency(p)
	return p
}

/**
 * @Description: 插入
 * @Author: Allen
 * @param {*core.TradeDetail} trade
 * @return {*}
 * @error: 
 */
func (p *tradeDetailReposiroty) Insert(trade *core.TradeDetail) int64 {
	result := p.db.Create(trade)
	if result.Error != nil {
		log.WithFields(log.Fields{
			"err": result.Error,
		}).Fatalf("failed to insert TradeDetail")
	}
	return result.RowsAffected
}

/**
 * @Description: 按条件查询
 * @Author: Allen
 * @param {core.TradeDetail} trade
 * @return {*}
 * @error: 
 */
func (p *tradeDetailReposiroty) SelectByCondition(trade core.TradeDetail) *core.TradeDetail {
	result := core.TradeDetail{}
	p.db.Where(&trade).First(&result)
	return &result
}

/**
 * @Description: 根据ID更新数据
 * @Author: Allen
 * @param {core.TradeDetail} trade
 * @return {*}
 * @error: 
 */
func (p *tradeDetailReposiroty) UpdateById(trade core.TradeDetail) int64 {
	result := p.db.Model(&trade).Where("id = ?", trade.ID).Updates(trade)
	if result.Error != nil {
		log.WithFields(log.Fields{
			"err": result.Error,
		}).Fatalf("failed to update TradeDetail")
	}
	return result.RowsAffected
}