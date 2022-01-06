/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-05-25 19:20:21
 * @LastEditors: Allen
 * @LastEditTime: 2021-05-25 19:20:21
 */
package infra

import log "github.com/sirupsen/logrus"

type DataBase interface {
	SetupDB() error
}


func SetupDB(database DataBase) error {
	err := database.SetupDB()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatalf("failed to setupDB")
		
		return err
	}
	return nil
}