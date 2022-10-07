package forker

import "errors"

var (
	ErrReuseportOnWindows = errors.New("please enable reuseport for windows")
	ErrOverRecovery       = errors.New("exceeding recovery child of forker")
)
