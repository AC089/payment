/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-05-27 17:02:16
 * @LastEditors: Allen
 * @LastEditTime: 2021-05-27 17:02:20
 */
package infra

import (
	"fmt"
	"payment/pkg/common/util"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Host string `yaml:"host"`
	Port uint32 `yaml:"port"`
	Password string `yaml:"password"`
	Db int `yaml:"db"`
	DialTimeout time.Duration `yaml:"dialTimeout"`
	ReadTimeout time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
	PoolSize int `yaml:"poolSize"`
	MinIdleConns int `yaml:"minIdleConns"`
	MaxConnAge time.Duration `yaml:"maxConnAge"`
	PoolTimeout time.Duration `yaml:"poolTimeout"`
}

func (r *RedisConfig) SetupRedis() {
	rdb := redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%d", r.Host, r.Port),
        Password: r.Password, // no password set
        DB:       r.Db,  // use default DB
    })

	if r.DialTimeout != 0 {
		rdb.Options().DialTimeout = r.DialTimeout
	}
	if r.ReadTimeout != 0 {
		rdb.Options().ReadTimeout = r.ReadTimeout
	}
	if r.WriteTimeout != 0 {
		rdb.Options().WriteTimeout = r.WriteTimeout
	}
	if r.PoolSize != 0 {
		rdb.Options().PoolSize = r.PoolSize
	}
	if r.MinIdleConns != 0 {
		rdb.Options().MinIdleConns = r.MinIdleConns
	}
	if r.MaxConnAge != 0 {
		rdb.Options().MaxConnAge = r.MaxConnAge
	}
	if r.PoolTimeout != 0 {
		rdb.Options().PoolTimeout = r.PoolTimeout
	}
	r.register(rdb)
}

func (r *RedisConfig) register(rdb *redis.Client) {
	util.GetRedis().SetClient(rdb)
}