package ljconf

import (
	"testing"

	"github.com/funkygao/assert"
)

func TestLoadFromZookeeper(t *testing.T) {
	ZkSvr = "localhost:2181"
	cf, err := Load("/_test_ljconf")
	assert.Equal(t, nil, err)
	assert.Equal(t, 58, cf.Int("hello", 0))

	assert.Equal(t, "/_test_ljconf", cf.ConfPath().S())
}
