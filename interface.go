package forker

import (
	"github.com/labstack/echo/v4"
	"net"
	"os/exec"
)

type Forker interface {
	ListenAndServe(address string) error
	ListenAndServeTLS(address, certFile, keyFile string) error
	NumOfChild() int
	ChildPids() []int

	forker(address string) (err error)
	listen(address string) (net.Listener, error)
	setTCPListenerFiles(address string) error
	doCmd() (*exec.Cmd, error)
}

type EchoForker interface {
	Start(address string) error
	GetEcho() *echo.Echo
}
