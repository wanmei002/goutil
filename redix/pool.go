package redix

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

var pl *redis.Pool

type redix struct {
	conn redis.Conn
}

func NewRedisConn() *redix {
	s := &redix{pl.Get()}

	return s
}

func (x *redix) Get(key string) string {
	v, err := redis.String(x.conn.Do("GET", key))
	if err != nil {
		log.Println("redis get failed; err:", err)
	}
	return v
}

func (x *redix) Set(key, val string, expire ...int) error {

	if len(expire) > 0 {
		_, err := x.conn.Do("SET", key, val, "EX", strconv.Itoa(expire[0]))
		return err
	} else {
		_, err := x.conn.Do("SET", key, val)
		return err
	}
}

// 删除key
func (x *redix) Del(key string) error {
	_, err := x.conn.Do("DEL", key)
	return err
}

// 对key 设置过期时间
func (x *redix) Expire(key string, expTime int) error {
	_, err := x.conn.Do("EXPIRE", key, expTime)
	return err
}

func (x *redix) Close() {
	x.conn.Close()
}

func Pool(host, port string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     100,
		IdleTimeout: 100 * time.Second,
		Dial: func() (redis.Conn, error) {
			log.Printf("%v:%v", host, port)
			c, err := redis.Dial("tcp", fmt.Sprintf("%v:%v", host, port))
			if err != nil {
				return nil, err
			}
			_, err = c.Do("AUTH", "123456")
			if err != nil {
				return nil, err
			}

			return c, nil
		},
	}
}

func InitPool(host, port string) {
	pl = Pool(host, port)
}
