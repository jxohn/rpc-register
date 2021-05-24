package consumer

import (
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/samuel/go-zookeeper/zk"
)

var (
	zkConn *zk.Conn
)

func Register(zkAddress []string, name string) (result []string, err error) {
	connect, events, err := zk.Connect(zkAddress, 3*time.Second)

	if err != nil {
		return nil, errors.Wrap(err, "connect to zk failed")
	}

	for {
		select {
		case event := <-events:
			if event.State == zk.StateConnected {
				goto NEXT
			}
		case <-time.After(3 * time.Second):
			return nil, errors.New("timeout to connect to zk")
		}
	}

NEXT:
	zkConn = connect
	exists, _, err := zkConn.Exists(name)
	if !exists {
		return nil, errors.Errorf("not available provider of [%s]", name)
	}

	children, _, err := zkConn.Children(name)
	log.Println(children)
	for i := range children {
		get, _, err := zkConn.Get(name + "/" + children[i])
		if err != nil {
			log.Println(err)
			continue
		}
		result = append(result, string(get))
	}

	return result, nil
}
