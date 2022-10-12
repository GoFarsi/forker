package forker

import (
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

type Forker interface {
	child
	ListenAndServe(address string) error
	ListenAndServeTLS(address, certFile, keyFile string) error
}

type GrpcForker interface {
	child
	ServeGrpc(address string) error
	GetGrpcServer() *grpc.Server
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
