package data

import (
	"github.com/garyburd/redigo/redis"
	"log"
)

type CacheRepository struct {
	DB	redis.Conn
	server	string
	authpw	string
}

func (r *CacheRepository) ConnectCache() bool {
	r.DisconnectCache()

	redis_conn, err := redis.Dial("tcp", r.server)
	if err != nil {
		log.Printf("*** Cannot contact redis server at %s. No caching!", r.server)
	} else if r.authpw != "" {
		if _, err = redis_conn.Do("AUTH", r.authpw); err != nil {
			redis_conn.Close()
//			redis_conn = nil
			log.Print("*** Redis authentication failure. No caching!")
		}
	}
	if redis_conn != nil {
		log.Printf("Redis connection to %s established", r.server)
		r.DB = redis_conn
		return true
	}
	return false
}

func (r *CacheRepository) DisconnectCache() {
	if r.DB != nil {
		r.DB.Close()
		r.DB = nil
	}
}

func (r *CacheRepository) GetCache(key string) ([]byte, error) {
	cval, err := r.DB.Do("GET", key)
	if err != nil {
		log.Printf("Redis error on GET: %v", err)
		return nil, err
	}
	if cval == nil {
		return nil, err
	}
	return cval.([]byte), err
}

func (r *CacheRepository) SetCache(key string, value []byte) (err error) {
	resp, err := r.DB.Do("SET", key, value)
	if err != nil {
		log.Printf("Redis error on SET: %v", err)
		return
	}
	if resp != "OK" {
		var other_err redis.Error = "Did not get OK from redis SET"
		log.Printf(string(other_err))
		err = other_err
	}
	return
}
