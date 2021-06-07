package provider

import (
	"sync"
)

type L struct {
	inner sync.Mutex
}

func (l *L) Lock4that() {
	l.inner.Lock()
	defer l.inner.Unlock()
}
