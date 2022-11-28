package main

import (
	"container/list"
	"sync"
)

type ProxyList struct {
	data *list.List
	lock *sync.Mutex
}

func (pl *ProxyList) GetAll() []*Proxy {
	pl.lock.Lock()
	defer pl.lock.Unlock()
	ret := make([]*Proxy, 0)
	for e := pl.data.Front(); e != nil; e = e.Next() {
		ret = append(ret, e.Value.(*Proxy))
	}
	return ret
}

func (pl *ProxyList) Add(p *Proxy) bool {
	pl.lock.Lock()
	defer pl.lock.Unlock()
	for e := pl.data.Front(); e != nil; e = e.Next() {
		if e.Value.(*Proxy).Id == p.Id {
			return false
		}
	}
	pl.data.PushBack(p)
	return true
}

func (pl *ProxyList) Remove(id int64) *Proxy {
	pl.lock.Lock()
	defer pl.lock.Unlock()
	for e := pl.data.Front(); e != nil; e = e.Next() {
		p := e.Value.(*Proxy)
		if p.Id == id {
			pl.data.Remove(e)
			return p
		}
	}
	return nil
}

func (pl *ProxyList) Len() int {
	pl.lock.Lock()
	defer pl.lock.Unlock()
	return pl.data.Len()
}

func (pl *ProxyList) Front() *Proxy {
	pl.lock.Lock()
	defer pl.lock.Unlock()
	return pl.data.Front().Value.(*Proxy)
}

func (pl *ProxyList) GetById(id int64) *Proxy {
	pl.lock.Lock()
	defer pl.lock.Unlock()
	for e := pl.data.Front(); e != nil; e = e.Next() {
		p := e.Value.(*Proxy)
		if p.Id == id {
			return p
		}
	}
	return nil
}

func NewProxyList() *ProxyList {
	return &ProxyList{
		data: list.New(),
		lock: &sync.Mutex{},
	}
}
