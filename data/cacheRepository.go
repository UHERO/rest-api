package data

import (
	"github.com/garyburd/redigo/redis"
	"log"
)

type CacheRepository struct {
	DB *redis.Conn
}

func (r *CacheRepository) GetCache(key string) (value []byte, err error) {
	value, err := r.DB.Do("GET", key)
	if err != nil {
		log.Printf("*** Connection failure to Redis!")
	}
	return value.([]byte), err
}

func (r *CacheRepository) SetCache(key string, value []byte) err error {
	resp, err := r.DB.Do("SET", key, value)
	if err != nil {
		log.Printf("*** Connection failure to Redis!")
		return
	}
	if resp != "OK" {
		var othererr redis.Error
		othererr = "Did not get OK from redis SET"
		log.Printf(othererr)
		return othererr
	}
}

