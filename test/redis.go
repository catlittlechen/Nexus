// Copyright 2017
// Author: catlittlechen@gmail.com

package main

import (
	"encoding/json"
	"fmt"
	"github.com/catlittlechen/Nexus"
	"github.com/garyburd/redigo/redis"
	"time"
)

type config struct {
	Addr     string `json:"addr"`
	DB       int    `json:"db"`
	PoolSize int    `json:"poolsize"`
	Idle     int    `json:"idle"`
}

type Redis struct {
	cfg  string
	pool *redis.Pool
}

func NewRedis(cfg string) (nexus.DB, error) {
	var obj config
	fmt.Println(cfg)
	err := json.Unmarshal([]byte(cfg), &obj)
	if nil != err {
		return nil, err
	}

	return &Redis{
		cfg:  cfg,
		pool: NewRedisPool(obj.Addr, obj.DB, obj.PoolSize, obj.PoolSize),
	}, nil
}

func (r *Redis) Set(key string, value string) error {
	conn := r.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	return err
}

func (r *Redis) Get(key string) (string, error) {
	conn := r.pool.Get()
	defer conn.Close()

	value, err := redis.String(conn.Do("GET", key))
	if err == redis.ErrNil {
		err = nexus.ErrKeyNotFound
	}
	return value, err
}

func (r *Redis) Del(key string) error {
	conn := r.pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}

func (r *Redis) Close() error {
	err := r.pool.Close()
	return err
}

func (r *Redis) String() string {
	return r.cfg
}

// NewRedisPool 初始化Redis链接池
func NewRedisPool(addr string, db, poolsize, idle int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     poolsize,
		IdleTimeout: time.Duration(idle) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				fmt.Printf("Redis Dial Error: %s\n", err.Error())
				return nil, err
			}

			if _, err = c.Do("SELECT", db); err != nil {
				fmt.Printf("Redis SELECT Error: %s\n", err.Error())
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
