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

func (r *CacheRepository) ConnectCache() bool {
	r.DisconnectCache()

	redis_conn, err := redis.Dial("tcp", r.server)
	if err != nil {
		log.Printf("*** Cannot contact redis server at %s. No caching!", r.server)
	} else if r.authpw != "" {
		if _, err = redis_conn.Do("AUTH", r.authpw); err != nil {
			redis_conn.Close()
			redis_conn = nil
			log.Print("*** Redis authentication failure. No caching!")
		}
	}
	r.DB = redis_conn
	if r.DB != nil {
		log.Printf("Redis connection to %s established", r.server)
		return true
	}
	return false
}

func (r *CacheRepository) GetCache(key string) ([]byte, error) {
	cval, err := r.DB.Do("GET", key)
	if err != nil {
		log.Printf("Redis error on GET: %v. Retrying.", err)
		if r.ConnectCache() {
			return r.GetCache(key)
		}
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
		log.Printf("Redis error on SET: %v. Retrying.", err)
		if r.ConnectCache() {
			return r.SetCache(key, value)
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
