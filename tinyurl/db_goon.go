// https://godoc.org/github.com/mjibson/goon
package main

import (
	"time"

	"net/http"
	"github.com/mjibson/goon"
	"google.golang.org/appengine/log"
	"strconv"
)

type Link struct {
	Id string `datastore:"-" goon:"id"`
	Content string `datastore:"Content,noindex"`
	Date    time.Time `datastore:"startTime"`
}

type linkDbGoon struct{
	goon.Goon goon
}

// var _ LinkDb = &linkDbCloud{}

func newDbCloud(r *http.Request) *linkDbCloud {
	g := goon.NewGoon(r)
	return &linkDbCloud{goon: g}
}

func (db *linkDbCloud) GetLink(key string, r *http.Request) (string, error) {
	g := goon.NewGoon(r)

	l := new(Link)
	l.Id = key

	err := g.Get(l)
	if err != nil {
		log.Infof(c, "%v", err)
		return "", err
	}
	return l.Content, nil
}

func (db *linkDbCloud) AddLink(l string, r *http.Request) (string, error) {
	g := goon.NewGoon(r)
	l := Link{
		Id: "aa",
		Content: l,
		Date:    time.Now(),
	}

	if _, err := g.Put(l); err != nil {
		log.Infof(c, "%v", err)
		return "", err
	}
	return l.Id, nil
}

func (db *linkDbCloud) Close() error {
	return nil
}
