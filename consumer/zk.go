package consumer

import (
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/pkg/errors"
	"github.com/samuel/go-zookeeper/zk"
)

var (
	zkConn *zk.Conn
	path   string
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
	path = name
	exists, _, err := zkConn.Exists(name)
	if !exists || err != nil {
		return nil, errors.Errorf("not available provider of [%s]", name)
	}

	children, _, err := zkConn.Children(name)
	if err != nil {
		return nil, errors.Wrap(err, "cant get children")
	}
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

func GetConn() (string, error) {
	children, _, err := zkConn.Children(path)
	if err != nil {
		return "", errors.Wrapf(err, "can't get child of [%s]", path)
	}

	result := make([]string, 0)
	for i := range children {
		get, _, err := zkConn.Get(path + "/" + children[i])
		if err != nil {
			log.Println(err)
			continue
		}
		result = append(result, string(get))
	}

	ticker := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-ticker.C:
			return "", errors.New("failed to get one")
		default:

		}
		intn := rand.Intn(len(result))
		s := result[intn]
		_, err := net.DialTimeout("tcp", s, 500*time.Millisecond)
		if err != nil {
			log.Println("cant get!!!")
			continue
		}
		return result[intn], nil
	}

}
