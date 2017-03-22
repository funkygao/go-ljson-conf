package ljconf

import (
	"bytes"
	"strings"
	"time"

	"github.com/daviddengcn/go-villa"
	"github.com/daviddengcn/ljson"
	"github.com/samuel/go-zookeeper/zk"
)

var (
	// ZkSvr is zookeeper connection string
	// If not empty, will load from zookeeper instead of the default file.
	ZkSvr string
)

func loadFromZk(znode string) (conf *Conf, err error) {
	conn, _, err := zk.Connect(strings.Split(ZkSvr, ","), time.Second*20)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	data, _, err := conn.Get(znode)
	if err != nil {
		return nil, err
	}

	conf = &Conf{
		path: villa.Path(znode),
		db:   make(map[string]interface{}),
	}

	dec := ljson.NewDecoder(newRcReader(bytes.NewReader(data)))
	err = dec.Decode(&conf.db)
	return conf, err
}
