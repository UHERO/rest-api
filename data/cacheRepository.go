package data

import (
	"github.com/garyburd/redigo/redis"
	"log"
)

type CacheRepository struct {
	DB	redis.Conn
	Server	string
	Authpw	string
}

func (r *CacheRepository) ConnectCache() (bool, error) {
	r.DisconnectCache()

	redis_conn, err := redis.Dial("tcp", r.Server)
	if err != nil {
		log.Printf("*** Cannot contact redis server at %s. No caching!", r.Server)
	} else if r.Authpw != "" {
		if _, err = redis_conn.Do("AUTH", r.Authpw); err != nil {
			redis_conn.Close()
			redis_conn = nil
			log.Print("*** Redis authentication failure. No caching!")
		}
	}
	r.DB = redis_conn
	if r.DB != nil {
		log.Printf("Redis connection to %s established", r.Server)
		return true, nil
	}
	return false, err
}

func (r *CacheRepository) GetCache(key string, iteration int) ([]byte, error) {
	if r.DB == nil {
		if connected, err := r.ConnectCache(); !connected {
			return nil, err
		}
	}
	cached_val, err := r.DB.Do("GET", key)
	if err != nil {
		log.Printf("Redis error on GET: %v", err)
		if iteration < 5 {
			log.Print("Retrying connection...")
			if connected, _ := r.ConnectCache(); connected {
				return r.GetCache(key, iteration + 1)
			}
		}
		return nil, err
	}
	if cached_val == nil {
		return nil, err
	}
	return cached_val.([]byte), err
}

func (r *CacheRepository) SetCache(key string, value []byte, iteration int) (err error) {
	resp, err := r.DB.Do("SET", key, value)
	if err != nil {
		log.Printf("Redis error on SET: %v", err)
		if iteration < 5 {
			log.Print("Retrying connection...")
			if connected, _ := r.ConnectCache(); connected {
				return r.SetCache(key, value, iteration + 1)
			}
		}
		return
	}
	if resp != "OK" {
		var other_err redis.Error = "Did not get OK from redis SET"
		log.Print(string(other_err))
		err = other_err
	}
	return
}

func (r *CacheRepository) DisconnectCache() {
	if r.DB != nil {
		r.DB.Close()
		r.DB = nil
	}
}
