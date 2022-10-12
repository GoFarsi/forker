package forker

import (
	"github.com/labstack/echo/v4"
)

type Forker interface {
	child
	ListenAndServe(address string) error
	ListenAndServeTLS(address, certFile, keyFile string) error
}

type EchoForker interface {
	child
	StartEcho(address string) error
	GetEcho() *echo.Echo
}

type child interface {
	NumOfChild() int
	ChildPids() []int
}
