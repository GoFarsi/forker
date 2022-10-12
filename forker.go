package forker

import (
	"flag"
	"github.com/labstack/echo/v4"
	"github.com/valyala/fasthttp/reuseport"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

var _ Forker = (*Fork)(nil)

const (
	childFlag       = "-child"
	_defaultNetwork = TCP4
)

type Fork struct {
	echo *echo.Echo

	serveFunc    func(ln net.Listener) error
	serveTLSFunc func(ln net.Listener, certFile, keyFile string) error

	recoverChild int
	ln           net.Listener
	files        []*os.File

	childsPid []int

	Network   Network // Network is net type tcp4, tcp, tcp6, udp, udp4, udp6
	ReusePort bool    // ReusePort use for windows support child process base on system call
}

type processSignal struct {
	pid int
	err error
}

func init() {
	flag.Bool(childFlag[1:], false, "is forker child process")
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

// NewEchoForker create echo forker object
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

func (f *Fork) forker(address string) (err error) {
	if !f.ReusePort {
		if runtime.GOOS == "windows" {
			return ErrReuseportOnWindows
		}

		if err = f.setTCPListenerFiles(address); err != nil {
			return
		}

		defer func() {
			if e := f.ln.Close(); e != nil {
				err = e
			}
		}()
	}

	goMaxProcess := runtime.GOMAXPROCS(0)
	sigCh := make(chan processSignal, goMaxProcess)
	childProcess := make(map[int]*exec.Cmd)

	defer func() {
		for _, proc := range childProcess {
			_ = proc.Process.Kill()
		}
	}()

	for i := 0; i < goMaxProcess; i++ {
		var cmd *exec.Cmd
		cmd, err = f.doCmd()
		if err != nil {
			return err
		}

		pid := cmd.Process.Pid
		childProcess[pid] = cmd
		f.childsPid = append(f.childsPid, pid)

		go func() {
			sigCh <- processSignal{cmd.Process.Pid, cmd.Wait()}
		}()
	}

	var exitedProcess int
	for sig := range sigCh {
		delete(childProcess, sig.pid)

		if exitedProcess++; exitedProcess > f.recoverChild {
			err = ErrOverRecovery
			break
		}

		var cmd *exec.Cmd
		if cmd, err = f.doCmd(); err != nil {
			break
		}

		childProcess[cmd.Process.Pid] = cmd
		go func() {
			sigCh <- processSignal{cmd.Process.Pid, cmd.Wait()}
		}()
	}

	return nil
}

func (f *Fork) listen(address string) (net.Listener, error) {
	runtime.GOMAXPROCS(1)

	if f.ReusePort {
		return reuseport.Listen(f.Network.String(), address)
	}

	return net.FileListener(os.NewFile(3, ""))
}

func (f *Fork) setTCPListenerFiles(address string) error {

	tcpAddr, err := net.ResolveTCPAddr(f.Network.String(), address)
	if err != nil {
		return err
	}

	tcpListener, err := net.ListenTCP(f.Network.String(), tcpAddr)
	if err != nil {
		return err
	}

	f.ln = tcpListener

	file, err := tcpListener.File()
	if err != nil {
		return err
	}

	f.files = []*os.File{file}

	return nil
}

func (f *Fork) doCmd() (*exec.Cmd, error) {
	cmd := exec.Command(os.Args[0], append(os.Args[1:], childFlag)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.ExtraFiles = f.files
	return cmd, cmd.Start()
}

func isChild() bool {
	for _, arg := range os.Args[1:] {
		if arg == childFlag {
			return true
		}
	}
	return false
}
