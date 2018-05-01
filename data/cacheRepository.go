package data

import (
	"errors"
	"log"

	"github.com/garyburd/redigo/redis"
)

type CacheRepository struct {
	Pool *redis.Pool
	TTL  int
}

func (r *CacheRepository) GetCache(key string) ([]byte, error) {
	c := r.Pool.Get()
	defer c.Close()
	value, err := c.Do("GET", key)
	if err != nil {
		log.Printf("Redis error on GET: %v", err)
		return nil, err
	}
	if value == nil {
		log.Printf("Redis cached val nil on GET: %v", err)
		return nil, err
	}
	log.Printf("Redis GET: %s", key)
	return value.([]byte), err
}

func (r *CacheRepository) SetCache(key string, value []byte) (err error) {
	c := r.Pool.Get()
	defer c.Close()
	c.Send("MULTI")
	c.Send("SET", key, value)
	c.Send("SET", key+":fresh", value)
	c.Send("EXPIRE", key+":fresh", r.TTL)
	response, err := redis.Values(c.Do("EXEC"))
	log.Print(response)
	if err != nil {
		log.Printf("Redis error on SET or EXPIRE: %v", err)
		return
	}
	var setResponse string
	var setResponseFresh string
	var expireResponse int
	if _, err := redis.Scan(response, &setResponse, &setResponseFresh, &expireResponse); err != nil {
		log.Print("Error on scan of redis response")
	}
	if setResponse != "OK" {
		err = errors.New("Did not get OK from Redis SET")
		log.Print(err)
		return
	}
	if setResponseFresh != "OK" {
		err = errors.New("Did not get OK from Redis SET")
		log.Print(err)
		return
	}
	if expireResponse != 1 {
		log.Printf("Did not set expiration to %v", r.TTL)
	}
	log.Printf("Redis SET: %s", key)
	return
}
