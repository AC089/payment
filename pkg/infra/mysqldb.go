/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-05-26 14:59:27
 * @LastEditors: Allen
 * @LastEditTime: 2021-05-26 14:59:31
 */
package infra

import (
	"fmt"
	"payment/pkg/infra/repository"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlDB struct {
	Host string `yaml:"host"`
	Port uint32 `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Dbname string `yaml:"db"`
	MaxIdle int `yaml:"maxIdle"`
	MaxOpen int `yaml:"maxOpen"`
	MaxLifetime time.Duration `yaml:"maxLifetime"`
}

func (m *MysqlDB) SetupDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", m.Username, m.Password, m.Host, m.Port, m.Dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
  
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatalf("failed to open mysql connection")
		return err
	}
	
	sqlDB, err := db.DB()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatalf("failed to create mysql connection pool")
		return err
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(m.MaxIdle)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(m.MaxOpen)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(m.MaxLifetime * time.Minute)
	//注册仓储
	m.register(db)
	return nil
}

func (m *MysqlDB) register(db *gorm.DB)  {
	repository.AliConfReposirotyInstance().SetDB(db).Inject()
	repository.WechatConfReposirotyInstance().SetDB(db).Inject()
	repository.PayAppReposirotyInstance().SetDB(db).Inject()
	// repository.GetReceiptDetailReposiroty().SetDB(db)
	repository.RefundDetailReposirotyInstance().SetDB(db).Inject()
	repository.TradeDetailReposirotyInstance().SetDB(db).Inject()
}  