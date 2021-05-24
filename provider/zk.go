package provider

import (
	"time"

	"github.com/pkg/errors"
	"github.com/samuel/go-zookeeper/zk"
)

var (
	zkConn *zk.Conn
)

func Online(zkAddress []string, name string, hostPort string) error {
	connect, events, err := zk.Connect(zkAddress, 3*time.Second)

	if err != nil {
		return errors.Wrap(err, "connect to zk failed")
	}

	for {
		select {
		case event := <-events:
			if event.State == zk.StateConnected {
				goto NEXT
			}
		case <-time.After(3 * time.Second):
			return errors.New("timeout to connect to zk")
		}
	}

NEXT:

	zkConn = connect
	exists, _, err := zkConn.Exists(name)
	if !exists {
		path, err := zkConn.Create(name, nil, 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			return errors.Wrapf(err, "failed to create zk path [%s]", name)
		}
		if path != name {
			return errors.Errorf("create [%s] not equals [%s]", path, name)
		}
	}

	_, err = zkConn.CreateProtectedEphemeralSequential(name+"/"+hostPort, []byte(hostPort), zk.WorldACL(zk.PermAll))
	if err != nil {
		return errors.Wrapf(err, "failed to create child path")
	}

	return nil
}

func Offline(name string, hostPort string) error {
	return nil
}
