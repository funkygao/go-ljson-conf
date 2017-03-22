package ljconf

type option func()

// WithZkSvr accepts svr as comma separated zookeeper address and use
// zookeeper as underlying configuration repo storage.
func WithZkSvr(svr string) option {
	return func() {
		zkSvr = svr
	}
}
