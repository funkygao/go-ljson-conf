package ljconf

var dsn string

func Load(path string, options ...option) (conf *Conf, err error) {
	for _, opt := range options {
		opt()
	}

	dsn = path

	if len(zkSvr) > 0 {
		return loadFromZk(path)
	}

	return loadFromFile(path)
}
