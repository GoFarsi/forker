package forker

import (
	"flag"
	"github.com/valyala/fasthttp/reuseport"
	"net"
	"os"
	"os/exec"
	"runtime"
)

const (
	childFlag       = "-child"
	_defaultNetwork = TCP4
)

type processSignal struct {
	pid int
	err error
}

func init() {
	flag.Bool(childFlag[1:], false, "is forker child process")
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
