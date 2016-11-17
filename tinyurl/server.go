package main

import (
	"net/http"
	"github.com/labstack/echo"
	"crypto/md5"
	"fmt"
	"github.com/labstack/gommon/log"
	"io"
	"html/template"
)

type (
	link struct {
		Key string `json:"key"`
		Link string `json:"link"`
	}
	Template struct {
		templates *template.Template
	}
)

var (
	links map[string]link
)

func GetLink(c echo.Context) error {
	key := c.Param("key")
	l,e := links[key]
	if e!=false {
		return c.String(http.StatusOK, l.Link)
	} else {
		return c.String(http.StatusNotFound, "not found")
	}
}

func CreateLink(c echo.Context) error {
	l := c.FormValue("link")
	key :=  fmt.Sprintf("%x", md5.Sum([]byte(l)))
	links[key] = link{
		Key: key,
		Link: l,
	}
	log.Printf("shorten_link:" + key)
	return c.String(http.StatusOK, "shorten_link:" + key) // TODO
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func top(c echo.Context) error {
	return c.Render(http.StatusOK, "top", "tinyurl")
}

func main() {
	links = make(map[string]link)

	e := echo.New()
	e.POST("/create", CreateLink)
	e.GET("/:key", GetLink)

	t := &Template{
		templates: template.Must(template.ParseFiles("tmp/index.html")),
	}
	e.Renderer = t
	e.GET("/", top)

	if err := e.Start(":9000"); err != nil {
		e.Logger.Fatal(err.Error())
	}
}
