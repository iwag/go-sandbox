package main

import (
	"html/template"
	"net/http"
	"time"

	"appengine"
	"appengine/datastore"
)

type Link struct {
	Id string
	Content string
	Date    time.Time
}

type linkDbCloud struct{
}

var _ LinkDb = &linkDbCloud{}

func NewClient(addr string) (LinkDb, error) {
	return &linkDbCloud{}, nil
}

func (db *linkDbCloud) GetLink(key string) (string, error) {
	c := appengine.NewContext()
	q := datastore.NewQuery("Links").Ancestor(linkKey(c)).Order("-Date").Limit(10)
	links := make([]Link, 0, 10)
	if _, err := q.GetAll(c, &links); err != nil {
		return "", err
	}
	return links[0].Content
}

func (db *linkDbCloud) AddLink(l string) (string, err error) {

	c := appengine.NewContext()
	key := datastore.NewIncompleteKey(c, "Links", linkKey(c))
	g := Link{
		Id: "aaa",
		Content: l,
		Date:    time.Now(),
	}
	_, err := datastore.Put(c, key, &g)
	if err != nil {
		return "", err
	}
	return g.Id, nil
}

func (db *linkDbCloud) Close() error {
	return nil
}


func linkKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Links", "default_guestbook", 0, nil)
}
