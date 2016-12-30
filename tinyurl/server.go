package main

import (
	"net/http"
	"gopkg.in/gin-gonic/gin.v1"
	"html/template"
	"log"
)

type (
	Template struct {
		templates *template.Template
	}
)

type LinkDb interface {
	GetLink(string, *http.Request) (string, error)

	AddLink(string, *http.Request) (string, error)

	Close() error
}

var (
	db LinkDb
)

//func init() {
//	db = newDbCloud() //newDbMem()
//}

func redirectToLink(c *gin.Context) {
	key := c.Param("key")

	if l,e := db.GetLink(key, c.Request); e == nil {
		c.Redirect(http.StatusMovedPermanently, l)
	} else {
		c.String(http.StatusNotFound, "not found")
	}
}

func createLink(c *gin.Context) {
	l := c.PostForm("link")
	key, err := db.AddLink(l, c.Request);
	if err != nil {
		c.String(http.StatusNoContent, "not found")
	} else {
		log.Printf("shorten_link:" + key)
		c.String(http.StatusOK, "http://localhost:8080/" + key) // TODO
	}
}

func top(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", struct { Name string } { Name: "tinyurl" } )
}

func init() {
	db = newDbGoon() // newDbCloud() // newDbMem()

	// if appengine
	// db = newDbMem()
	// r := gin.New()
	r := gin.Default()

	r.LoadHTMLGlob("tmp/*")

	r.GET("/", top)
	r.POST("/create", createLink)
	r.GET("/:key", redirectToLink)

	http.Handle("/", r)
	//r.Run(":8080")
}
