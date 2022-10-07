package forker

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"testing"
)

func Test_Forker(t *testing.T) {
	t.Parallel()
	setup()
	defer tearDown()

	srv := &http.Server{}
	f := New(srv, WithReusePort(true))

	addr := fmt.Sprintf("0.0.0.0:%d", rand.Intn(9000-3000)+3000)

	err := f.ListenAndServe(addr)
	if err != nil {
		t.Error(err)
	}

	f.(*Fork).ln.Close()

	lnAddr := f.(*Fork).ln.Addr().String()
	if lnAddr != addr {
		t.Errorf("want forker address %s, but listener address is %s", addr, lnAddr)
	}

	if f.(*Fork).ln == nil {
		t.Error("listener is null")
	}
}

func setup() {
	os.Args = append(os.Args, childFlag)
}

func tearDown() {
	os.Args = os.Args[:len(os.Args)-1]
}
