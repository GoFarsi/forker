package main

import (
	"github.com/Ja7ad/forker"
	"github.com/labstack/echo/v4"
)

func main() {
	f := forker.NewEchoForker()
	e := f.GetEcho()

	e.GET("/", Greeting)

	e.Logger.Fatal(f.Start(":8080"))
}

func Greeting(c echo.Context) error {
	return c.String(200, "greeting")
}
