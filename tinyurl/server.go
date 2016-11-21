package main

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"io"
	"html/template"
)

type (
	Template struct {
		templates *template.Template
	}
)

type LinkDb interface {
	GetLink(key string) (string, error)

	AddLink(string) (string, error)

	Close() error
}

var db LinkDb = linkDbMem{
	links: make(map[string]string),
}

func redirectToLink(c echo.Context) error {
	key := c.Param("key")

	if l,e := db.GetLink(key); e!=nil {
		return c.Redirect(http.StatusMovedPermanently, l)
	} else {
		return c.String(http.StatusNotFound, "not found")
	}
}

func createLink(c echo.Context) error {
	l := c.FormValue("link")
	if 	key, err := db.AddLink(l); err != nil {
		log.Errorf("shorten_link: %s", key)
		return c.NoContent(http.StatusNoContent)
	} else {
		log.Printf("shorten_link:" + key)
		return c.String(http.StatusOK, "http://localhost:9000/" + key) // TODO
	}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func top(c echo.Context) error {
	return c.Render(http.StatusOK, "top", "tinyurl")
}

func main() {
	e := echo.New()
	e.POST("/create", createLink)
	e.GET("/:key", redirectToLink)

	t := &Template{
		templates: template.Must(template.ParseFiles("tmp/index.html")),
	}
	e.Renderer = t
	e.GET("/", top)

	if err := e.Start(":9000"); err != nil {
		e.Logger.Fatal(err.Error())
	}
}
