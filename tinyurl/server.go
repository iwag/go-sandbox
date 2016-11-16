package main

import (
	"net/http"
	"github.com/labstack/echo"
)

func getLink(c echo.Context) error {
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}

func createLink(c echo.Context) error {
	id := c.FormValue("link")
	return c.String(http.StatusOK, "unique_link:" + id) // TODO
}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to tinyurl")
	})
	e.POST("/create", createLink)
	e.GET("/:id", getLink)

	if err := e.Start(":9000"); err != nil {
		e.Logger.Fatal(err.Error())
	}
}
