package data

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"errors"
)

type CacheRepository struct {
	Pool	*redis.Pool
}

func (r *CacheRepository) GetCache(key string, iteration int) ([]byte, error) {
	conn := r.Pool.Get()
	defer conn.Close()
	cached_val, err := conn.Do("GET", key)
	if err != nil {
		log.Printf("Redis error on GET: %v", err)
		return nil, err
	}
	if cached_val == nil {
		log.Printf("Redis cached val nil on GET: %v", err)
		return nil, err
	}
	log.Printf("Redis GET: %s", key)
	return cached_val.([]byte), err
}

func (r *CacheRepository) SetCache(key string, value []byte, iteration int) (err error) {
	conn := r.Pool.Get()
	defer conn.Close()
	resp, err := conn.Do("SET", key, value)
	if err != nil {
		log.Printf("Redis error on SET: %v", err)
		return
	}
	if resp != "OK" {
		err = errors.New("Did not get OK from Redis SET")
		log.Print(err)
		return
	}
	log.Printf("Redis SET: %s", key)
	return
}
