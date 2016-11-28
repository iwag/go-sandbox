// +build !appengine

package main

import (
	"gopkg.in/redis.v5"
	"fmt"
)

type linkDbRedis struct{
	Client redis.Client
}

var _ LinkDb = &linkDbRedis{}

func NewClient(addr string) (LinkDb, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	return &linkDbRedis{
		Client: client,
	},nil
}

func (db *linkDbRedis) GetLink(key string) (string, error) {
	l, err := db.Client.Get(key).Result()
	if err != nil {
		return 0, fmt.Errorf("%v", err)
	}

	return l
}

func (db *linkDbRedis) AddLink(l string) (string, err error) {
	key, err := db.Client.RandomKey().Result()
	if err != nil {
		return 0, fmt.Errorf("%v", err)
	}

	if err := db.Client.Set(key, l).Err(); err != nil {
		return 0, fmt.Errorf("%v", err)
	}
	return key, nil
}

func (db *linkDbRedis) Close() error {
	return db.Client.Close()
}
