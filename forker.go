package forker

import (
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"runtime"
)

var _ Forker = (*Fork)(nil)

type Fork struct {
	echo *echo.Echo
	grpc *grpc.Server

	serveFunc    func(ln net.Listener) error
	serveTLSFunc func(ln net.Listener, certFile, keyFile string) error

	recoverChild int
	ln           net.Listener
	files        []*os.File

	childsPid []int

	Network   Network // Network is net type tcp4, tcp, tcp6, udp, udp4, udp6
	ReusePort bool    // ReusePort use for windows support child process base on system call
}

// New create forker for listen and serve http server
func New(httpServer *http.Server, opts ...Option) Forker {
	forker := &Fork{
		serveFunc:    httpServer.Serve,
		serveTLSFunc: httpServer.ServeTLS,
		recoverChild: runtime.GOMAXPROCS(0) / 2,
	}

	for _, opt := range opts {
		opt(forker)
	}

	if forker.Network == 0 {
		forker.Network = _defaultNetwork
	}

	return forker
}

// NewEchoForker create for listen and serve echo
func NewEchoForker(opts ...Option) EchoForker {
	forker := &Fork{
		echo:         echo.New(),
		recoverChild: runtime.GOMAXPROCS(0) / 2,
	}

	for _, opt := range opts {
		opt(forker)
	}

	if forker.Network == 0 {
		forker.Network = _defaultNetwork
	}

	return forker
}

// NewGrpcForker create for listen and serve grpc server
func NewGrpcForker(grpcServerOpts []grpc.ServerOption, opts ...Option) GrpcForker {
	forker := &Fork{
		recoverChild: runtime.GOMAXPROCS(0) / 2,
	}

	for _, opt := range opts {
		opt(forker)
	}

	grpcOpts := make([]grpc.ServerOption, 0)

	if len(grpcServerOpts) != 0 {
		grpcOpts = append(grpcOpts, grpcServerOpts...)
	}

	forker.grpc = grpc.NewServer(grpcOpts...)

	if forker.Network == 0 {
		forker.Network = _defaultNetwork
	}

	return forker
}

// ServeGrpc serve grpc server
func (f *Fork) ServeGrpc(address string) error {
	if isChild() {
		ln, err := f.listen(address)
		if err != nil {
			return err
		}
		f.ln = ln
		return f.grpc.Serve(ln)
	}
	return f.forker(address)
}

// GetGrpcServer return grpc server object
func (f *Fork) GetGrpcServer() *grpc.Server {
	return f.grpc
}

// StartEcho listener echo
func (f *Fork) StartEcho(address string) error {
	if isChild() {
		ln, err := f.listen(address)
		if err != nil {
			return err
		}
		f.echo.Listener = ln
		f.ln = ln
		return f.echo.Start(address)
	}
	return f.forker(address)
}

// GetEcho return echo object
func (f *Fork) GetEcho() *echo.Echo {
	return f.echo
}

// ListenAndServe listen and serve http server
func (f *Fork) ListenAndServe(address string) error {
	if isChild() {
		ln, err := f.listen(address)
		if err != nil {
			return err
		}
		f.ln = ln
		return f.serveFunc(ln)
	}
	return f.forker(address)
}

// ListenAndServeTLS listen and serve http server with tls support
func (f *Fork) ListenAndServeTLS(address, certFile, keyFile string) error {
	if isChild() {
		ln, err := f.listen(address)
		if err != nil {
			return err
		}

		f.ln = ln

		return f.serveTLSFunc(ln, certFile, keyFile)
	}

	return f.forker(address)
}

// NumOfChild number of child process
func (f *Fork) NumOfChild() int {
	return len(f.childsPid)
}

// ChildPids list child process PID
func (f *Fork) ChildPids() []int {
	return f.childsPid
}
