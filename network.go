package forker

//go:generate stringer -type=Network

type Network int

const (
	TCP4 Network = iota + 1
	TCP6
	UDP
	UDP4
	UDP6
)
