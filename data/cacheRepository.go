package data

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"log"
)

type CacheRepository struct {
	Pool 	*redis.Pool
	TTL  	int
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
		//log.Printf("Redis cached val nil on GET: %v", err)
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
	c.Send("EXPIRE", key, r.TTL)
	response, err := redis.Values(c.Do("EXEC"))
	if err != nil {
		log.Printf("Redis error on SET or EXPIRE: %v", err)
		return
	}
	var setResponse string
	var expireResponse int
	if _, err := redis.Scan(response, &setResponse, &expireResponse); err != nil {
		log.Print("Error on scan of redis response")
	}
	if setResponse != "OK" {
		err = errors.New("did not get OK from Redis SET")
		log.Print(err)
		return
	}
	if expireResponse != 1 {
		log.Printf("Did not set expiration to %v", r.TTL)
	}
	log.Printf("Redis SET: %s", key)
	return
}

// SetCachePair sets a pair of cache keys key (w/o expiration) and key:fresh (w/ expiration). Used for caching responses from the census proxy.
func (r *CacheRepository) SetCachePair(key string, value []byte) (err error) {
	r.SetCache(key+":fresh", value)
	r.SetCache(key, value)
	c := r.Pool.Get()
	defer c.Close()
	c.Send("PERSIST", key)
	return
}
