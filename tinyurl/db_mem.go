package main

import (
	"crypto/md5"
	"fmt"
)


type linkDbMem struct{
	links map[string]string
}

var _ LinkDb = &linkDbMem{}

func (db *linkDbMem) GetLink(key string) (string, error) {
	l, err := db.links[key]
	if err != false {
		return "", nil // TODO !!
	}

	return l, nil
}

func (db *linkDbMem) AddLink(l string) (string, error) {
	key :=  fmt.Sprintf("%x", md5.Sum([]byte(l)))
	db.links[key] = l
	return key, nil
}

func (db *linkDbMem) Close() error {
	return nil
}
