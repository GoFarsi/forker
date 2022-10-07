package forker

type Option func(f *Fork)

// WithReusePort enable reuse port option for windows
func WithReusePort(reusePort bool) Option {
	return func(f *Fork) {
		f.ReusePort = reusePort
	}
}

// WithCustomNetwork set network type listing type
func WithCustomNetwork(network Network) Option {
	return func(f *Fork) {
		f.Network = network
	}
}
