package main

import (
	"net/http"
	"github.com/labstack/echo"
	"io"
	"html/template"
	"os"
)

type (
	Template struct {
		templates *template.Template
	}
)

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func hello(c echo.Context) error {
	return c.Render(http.StatusOK, "hello", "world")
}

func main() {

	port := os.Getenv("PORT")

	e := echo.New()

	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = t
	e.GET("/", hello)

	if err := e.Start(":" + port); err != nil {
		e.Logger.Fatal(err.Error())
	}
}
