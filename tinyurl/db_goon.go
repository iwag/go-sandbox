// https://godoc.org/github.com/mjibson/goon
package main

import (
	"time"

	"net/http"
	"github.com/mjibson/goon"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine"
)

type Link struct {
	Id      string `datastore:"-" goon:"id"`
        Content string `datastore:"Content,noindex"`
        Date    time.Time `datastore:"startTime"`
}

type linkDbGoon struct {
	goon string
}

// var _ LinkDb = &linkDbCloud{}

func newDbGoon() *linkDbGoon {
	return &linkDbGoon{goon: ""}
}

func (db *linkDbGoon) GetLink(key string, r *http.Request) (string, error) {
	g := goon.NewGoon(r)
	c := appengine.NewContext(r)

	l := new(Link)
	l.Id = key

	err := g.Get(l)
	if err != nil {
		log.Infof(c, "%v", err)
		return "", err
	}
	return l.Content, nil
}

func (db *linkDbGoon) AddLink(l string, r *http.Request) (string, error) {
	g := goon.NewGoon(r)
	c := appengine.NewContext(r)
	link := Link{
		Id: "aa",
		Content: l,
		Date:    time.Now(),
	}

	if _, err := g.Put(&link); err != nil {
		log.Infof(c, "%v", err)
		return "", err
	}
	return link.Id, nil
}

func (db *linkDbGoon) Close() error {
	return nil
}
