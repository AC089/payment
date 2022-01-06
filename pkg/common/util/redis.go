/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-05-27 19:45:57
 * @LastEditors: Allen
 * @LastEditTime: 2021-05-27 19:46:00
 */
package util

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)


var(
	instance = new(_redis)
	ctx = context.Background()
)

type _redis struct {
	client *redis.Client
}

func GetRedis() *_redis {
	return instance
}

func (r *_redis) SetClient(client *redis.Client) {
	r.client = client
}

func (r *_redis) Get(key string) (string, error) {
	stat := r.client.Get(ctx, key)
	s, err := stat.Result()
	if err != nil {
		return "", err
	}
	return s, nil
}

func (r *_redis) Set(key string, value interface{}, expiration time.Duration) error {
	stat := r.client.Set(ctx, key, value, expiration)
	if stat.Err() != nil {
		return stat.Err()
	}
	t, err := stat.Result()
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"result": t,
	}).Debug("redis set result:")
	return nil
}

func (r *_redis) HMSet(key string, fields map[string]interface{}) error {
	stat := r.client.HMSet(ctx, key, fields)
	s, err := stat.Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}
	log.WithFields(log.Fields{
		"result": s,
		"err":    err,
	}).Debug("redis result:")

	return err
}

func (r *_redis) HExists(key string, field string) (bool, error) {
	stat := r.client.HExists(ctx, key, field)
	err := stat.Err()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	return stat.Val(), nil
}

func (r *_redis) HMGet(key string, fields ...string) ([]interface{}, error) {
	NIL := make([]interface{}, 0)
	stat := r.client.HMGet(ctx, key, fields...)
	s, err := stat.Result()
	if err != nil {
		if err == redis.Nil {
			return NIL, nil
		}
		return NIL, err
	}
	return s, nil
}

func (r *_redis) HGet(key string, fields string) (string, error) {
	stat := r.client.HGet(ctx, key, fields)
	s, err := stat.Result()
	if err != nil {
		return "", err
	}
	log.WithFields(log.Fields{
		"result": s,
	}).Debug("redis HGet result:")
	return s, nil
}

func (r *_redis) HSet(key string, field string, value string) error {
	stat := r.client.HSet(ctx, key, field, value)
	s, err := stat.Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}
	log.WithFields(log.Fields{
		"result": s,
	}).Debug("redis HGet result:")
	return nil
}

func (r *_redis) HDel(key string, fields ...string) error {
	stat := r.client.HDel(ctx, key, fields...)
	str, err := stat.Result()
	if err != redis.Nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}
	log.WithFields(log.Fields{
		"result": str,
	}).Debug("redis HDel result:")
	return nil
}

func (r *_redis) Del(key string) error {
	stat := r.client.Del(ctx, key)
	str, err := stat.Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}
	log.WithFields(log.Fields{
		"result": str,
	}).Debug("redis Del result:")
	return nil
}

func (r *_redis) LPush(key string, values ...interface{}) error {
	stat := r.client.LPush(ctx, key, values...)
	str, err := stat.Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}
	log.WithFields(log.Fields{
		"result": str,
	}).Debug("redis LPush result:")
	return nil
}

func (r *_redis) LPushX(key string, values ...interface{}) error {
	stat := r.client.LPushX(ctx, key, values...)
	str, err := stat.Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}
	log.WithFields(log.Fields{
		"result": str,
	}).Debug("redis LPushX result:")
	return nil
}

func (r *_redis) LRange(key string, start, stop int64) ([]string, error) {
	stat := r.client.LRange(ctx, key, start, stop)
	str, err := stat.Result()
	if err != nil {
		return nil, err
	}
	log.WithFields(log.Fields{
		"result": str,
	}).Debug("redis LRange result:")
	return str, nil
}

func (r *_redis) LPop(key string) (string, error) {
	stat := r.client.LPop(ctx, key)
	str, err := stat.Result()
	if err != nil {
		return "", err
	}
	log.WithFields(log.Fields{
		"result": str,
	}).Debug("redis LPop result:")
	return str, nil
}

func (r *_redis) RPop(key string) (string, error) {
	stat := r.client.RPop(ctx, key)
	str, err := stat.Result()
	if err != nil {
		return "", err
	}
	log.WithFields(log.Fields{
		"result": str,
	}).Debug("redis RPop result:")
	return str, nil
}

func (r *_redis) RPush(key string, values ...interface{}) error {
	stat := r.client.RPush(ctx, key, values...)
	str, err := stat.Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}
	log.WithFields(log.Fields{
		"result": str,
	}).Debug("redis RPush result:")
	return nil
}

func (r *_redis) LTrim(key string, start, stop int64) (string, error) {
	stat := r.client.LTrim(ctx, key, start, stop)
	str, err := stat.Result()
	if err != nil {
		return "", err
	}
	log.WithFields(log.Fields{
		"result": str,
	}).Debug("redis RPop result:")
	return str, nil
}

func (r *_redis) HInc(key string, fields string, value int) (int64, error) {
	stat := r.client.HIncrBy(ctx, key, fields, int64(value))
	str, err := stat.Result()
	if err != nil {
		return 0, err
	}
	log.WithFields(log.Fields{
		"result": str,
	}).Debug("redis HInc result:")
	return str, nil
}

func (r *_redis) Expire(key string, expiration time.Duration) error {
	stat := r.client.Expire(ctx, key, expiration)
	str, err := stat.Result()
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"result": str,
	}).Debug("redis Expire result:")
	return nil
}

func (r *_redis) Exists(key ...string) (int64, error) {
	stat := r.client.Exists(ctx, key...)
	num, err := stat.Result()
	if err != nil {
		return 0, err
	}
	log.WithFields(log.Fields{
		"result": num,
	}).Debug("redis Exists result:")
	return num, nil
}

func (r *_redis) ZAdd(key string, members ...*redis.Z) (int64, error) {
	stat := r.client.ZAdd(ctx, key, members...)
	num, err := stat.Result()
	if err != nil {
		return 0, err
	}
	log.WithFields(log.Fields{
		"result": num,
	}).Debug("redis ZAdd result:")
	return num, nil
}

func (r *_redis) ZRevRank(key, member string) (int64, error) {
	stat := r.client.ZRevRank(ctx, key, member)
	num, err := stat.Result()
	if err != nil {
		return 0, err
	}
	log.WithFields(log.Fields{
		"result": num,
	}).Debug("redis ZRevRank result:")
	return num, nil
}

func (r *_redis) ZScore(key, member string) (float64, error) {
	stat := r.client.ZScore(ctx, key, member)
	num, err := stat.Result()
	if err != nil {
		return 0, err
	}
	log.WithFields(log.Fields{
		"result": num,
	}).Debug("redis ZRevRank result:")
	return num, nil
}

func (r *_redis) ZCard(key string) (int64, error) {
	stat := r.client.ZCard(ctx, key)
	num, err := stat.Result()
	if err != nil {
		return 0, err
	}
	log.WithFields(log.Fields{
		"result": num,
	}).Debug("redis ZRevRank result:")
	return num, nil
}

func (r *_redis) ZRevRangeWithScores(key string, start, stop int64) ([]redis.Z, error) {
	NIL := make([]redis.Z, 0)
	stat := r.client.ZRevRangeWithScores(ctx, key, start, stop)
	sli, err := stat.Result()
	if err != nil {
		return NIL, err
	}
	log.WithFields(log.Fields{
		"result": sli,
	}).Debug("redis ZRevRangeWithScores result:")
	return sli, nil
}

func (r *_redis) Eval(script string, keys []string, args ...interface{}) (interface{}, error) {
	stat := r.client.Eval(ctx, script, keys, args...)
	tmp, err := stat.Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	log.WithFields(log.Fields{
		"result": tmp,
	}).Debug("redis Eval result:")
	return tmp, nil
}

func (r *_redis) Keys(pattern string) ([]string, error) {
	NIL := make([]string, 0)
	stat := r.client.Keys(ctx, pattern)
	sli, err := stat.Result()
	if err != nil {
		return NIL, err
	}
	return sli, nil
}

func (r *_redis) HKeys(key string) ([]string, error) {
	NIL := make([]string, 0)
	stat := r.client.HKeys(ctx, key)
	sli, err := stat.Result()
	if err != nil {
		return NIL, err
	}
	return sli, nil
}

func (r *_redis) HLen(key string) (int64, error) {
	stat := r.client.HLen(ctx, key)
	lens, err := stat.Result()
	if err != nil {
		return 0, err
	}
	return lens, nil
}

func (r *_redis) ScriptLoad(script string) (string, error) {
	stat := r.client.ScriptLoad(ctx, script)
	sha, err := stat.Result()
	if err != nil {
		return "", err
	}
	return sha, nil
}

func (r *_redis) EvalSha(sha string, keys []string, args ...interface{}) (interface{}, error) {
	stat := r.client.EvalSha(ctx, sha, keys, args...)
	tmp, err := stat.Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	log.WithFields(log.Fields{
		"result": tmp,
	}).Debug("redis Eval result:")
	return tmp, nil
}

func (r *_redis) Rename(key, newkey string) (int64, error) {
	stat := r.client.Rename(ctx, key, newkey)
	_, err := stat.Result()
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func (r *_redis) Persist(key string) error {
	stat := r.client.Persist(ctx, key)
	str, err := stat.Result()
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"result": str,
	}).Debug("redis Persist result:")
	return nil
}

func (r *_redis) Incr(key string) (int, error) {
	stat := r.client.Incr(ctx, key)
	rt, err := stat.Result()
	if err != nil {
		return 0, err
	}
	log.WithFields(log.Fields{
		"result": rt,
	}).Debug("redis Incr result:")
	return int(rt), nil
}

func (r *_redis) IncrBy(key string, value int64) error {
	stat := r.client.IncrBy(ctx, key, value)
	str, err := stat.Result()
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"result": str,
	}).Debug("redis Persist result:")
	return nil
}

/**
 * @Description: 批量查询hash表数据
 * @Author: Allen
 * @param {map[string]string} params
 * @return {*}
 * @error:
 */
func (r *_redis) HgetBatch(params map[string]string) (result map[string]map[string]string, err error) {
	result = make(map[string]map[string]string)
	pipe := r.client.Pipeline()
	tempMap1 := make(map[string]map[string]*redis.StringCmd)
	tempMap2 := make(map[string]*redis.StringCmd)
	for k, v := range params {
		tempMap2[v] = pipe.HGet(ctx, k, v)
		tempMap1[k] = tempMap2
	}
	_, err = pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return
	}

	for k1, v1 := range tempMap1 {
		for k2, v2 := range v1 {
			r, err := v2.Result()
			if err != nil && err != redis.Nil {
				return nil, err
			}
			if err != redis.Nil {
				result[k1][k2] = r
			}
		}
	}
	return result, nil
}