package core

import (
	"errors"
	"github.com/ZYallers/zgin/libraries/mvcs"
	"github.com/go-redis/redis"
	jsontime "github.com/liamylian/jsontime/v2/v2"
	"math/rand"
	"src/config/env"
	"src/libraries/util/helper"
	"strings"
	"time"
)

type Redis struct {
	mvcs.Redis
}

const hashAllFieldKey = "all"

var cache, session mvcs.RdsCollector

func (r *Redis) GetCache() *redis.Client {
	return r.NewClient(&cache, &env.Redis.Cache)
}

func (r *Redis) GetSession() *redis.Client {
	return r.NewClient(&session, &env.Redis.Session)
}

// 数据不存在情况下，为防止缓存雪崩，随机返回一个30到60秒的有效时间
func (r *Redis) NoDataExpiration() time.Duration {
	// 将时间戳设置成种子数
	rand.Seed(time.Now().UnixNano())
	return time.Duration(30+rand.Intn(30)) * time.Second
}

// 从String类型的缓存中读取数据，如没则重新调用指定方法重新从数据库中读取并写入缓存
func (r *Redis) CacheWithString(key string, output interface{}, fn func() (interface{}, bool), options ...interface{}) error {
	if val := r.GetCache().Get(key).Val(); val != "" {
		return jsontime.ConfigWithCustomTimeFormat.Unmarshal(helper.String2Bytes(val), &output)
	}

	var (
		isNull     bool
		data       interface{}
		expiration = env.Redis.CommonExpiration
	)

	if data, isNull = fn(); isNull {
		expiration = r.NoDataExpiration()
	} else {
		if len(options) > 0 {
			expiration = options[0].(time.Duration)
		}
	}

	var value string
	bte, err := jsontime.ConfigWithCustomTimeFormat.Marshal(data)
	if err != nil {
		value = "null"
	} else {
		value = helper.Bytes2String(bte)
		_ = jsontime.ConfigWithCustomTimeFormat.Unmarshal(bte, &output)
	}
	return r.GetCache().Set(key, value, expiration).Err()
}

//  DeleteCache 根据key删除对应缓存
//  @receiver r *Redis
//  @author Cloud|2021-12-07 13:56:08
//  @param key ...string ...
//  @return int64 ...
//  @return error ...
func (r *Redis) DeleteCache(key ...string) (int64, error) {
	return r.GetCache().Del(key...).Result()
}

//  HashGetAll ...
//  @receiver r *Redis
//  @author Cloud|2021-12-15 10:10:46
//  @param key string ...
//  @return result []interface{} ...
func (r *Redis) HashGetAll(key string) (result []interface{}) {
	all := r.GetCache().HGet(key, hashAllFieldKey).Val()
	if all == "" {
		return
	}
	keys := helper.RemoveDuplicateWithString(strings.Split(all, ","))
	if len(keys) == 0 {
		return
	}
	result = r.GetCache().HMGet(key, keys...).Val()
	return
}

//  HashMultiSet ...
//  @receiver r *Redis
//  @author Cloud|2021-12-15 10:10:49
//  @param key string ...
//  @param data map[string]interface{} ...
//  @return error ...
func (r *Redis) HashMultiSet(key string, data map[string]interface{}) error {
	fields := make([]string, 0)
	fieldValues := make(map[string]interface{}, 0)
	for k, v := range data {
		if k == "" || v == nil {
			continue
		}
		if b, err := jsontime.ConfigWithCustomTimeFormat.Marshal(v); err == nil {
			fieldValues[k] = helper.Bytes2String(b)
			fields = append(fields, k)
		}
	}

	if len(fields) == 0 {
		return errors.New("the data that can be saved is empty")
	}

	if val := r.GetCache().HGet(key, hashAllFieldKey).Val(); val != "" {
		fields = append(fields, strings.Split(val, ",")...)
	}

	var allFieldValue string
	if len(fields) > 0 {
		allFieldValue = strings.Join(helper.RemoveDuplicateWithString(fields), ",")
	}
	fieldValues[hashAllFieldKey] = allFieldValue
	return r.GetCache().HMSet(key, fieldValues).Err()
}

//  HashMultiDelete ...
//  @receiver r *Redis
//  @author Cloud|2021-12-15 10:49:15
//  @param key string ...
//  @param fields ...string ...
//  @return error ...
func (r *Redis) HashMultiDelete(key string, fields ...string) error {
	newFields := make([]string, 0)
	if val := r.GetCache().HGet(key, hashAllFieldKey).Val(); val != "" {
		newFields = append(newFields, strings.Split(val, ",")...)
	}
	if len(newFields) > 0 {
		for _, field := range fields {
			newFields = helper.RemoveWithString(newFields, field)
		}
	}

	var allFieldValue string
	if len(newFields) > 0 {
		allFieldValue = strings.Join(newFields, ",")
	}

	pl := r.GetCache().Pipeline()
	pl.HDel(key, fields...)
	pl.HSet(key, hashAllFieldKey, allFieldValue)
	_, err := pl.Exec()
	return err
}
