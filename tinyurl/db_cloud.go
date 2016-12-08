package main

import (
	"time"

	"net/http"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"strconv"
)

type Link struct {
	Content string
	Date    time.Time
}

type linkDbCloud struct{
	client string
}

// var _ LinkDb = &linkDbCloud{}

func newDbCloud() *linkDbCloud {
	return &linkDbCloud{client: ""}
}

func (db *linkDbCloud) GetLink(key string, r *http.Request) (string, error) {
	c := appengine.NewContext(r)
	i, e := strconv.ParseInt(key, 10, 0)
	if e != nil {
		log.Infof(c, "%v", e)
		return "", e
	}
	k := datastore.NewKey(c, "Links", "", i, nil)

	l := new(Link)
	err := datastore.Get(c, k, l)
	if err != nil {
		log.Infof(c, "%v", err)
		return "", err
	}
	return l.Content, nil
}

func (db *linkDbCloud) AddLink(l string, r *http.Request) (string, error) {
	c := appengine.NewContext(r)
	key := datastore.NewIncompleteKey(c, "Links", nil)
	g := Link{
		Content: l,
		Date:    time.Now(),
	}
	x, err := datastore.Put(c, key, &g)

	if err != nil {
		log.Infof(c, "%v", err)
		return "", err
	}
	return strconv.FormatInt(x.IntID(), 10), nil
}

func (db *linkDbCloud) Close() error {
	return nil
}
