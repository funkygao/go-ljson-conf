package ljconf

func Load(path string) (conf *Conf, err error) {
	if len(ZkSvr) > 0 {
		return loadFromZk(path)
	}

	return loadFromFile(path)
}
