package main

import (
	"context"
	"log"
	"time"

	"github.com/go-cmd/cmd"
)

const (
	DELAY_BEFORE_RUN_AGAIN = 5 * time.Second
	RUN_CONFIRM_DELAY      = 3 * time.Second
)

type ProxyRunStatusFunc func(id int64, runStatus string)
type LogsFunc func(id int64, time int64, _type string, message string)

type ProxyHandler struct {
	proxyList      *ProxyList
	clientFilePath string
	runStatusFunc  ProxyRunStatusFunc
	logsFunc       LogsFunc
	ShutDown       bool
	logs           map[int64]*LogsList
	// logsLock       *sync.RWMutex
}

func (ph *ProxyHandler) AddProxy(p *Proxy) {
	ph.RemoveProxy(p.Id)
	if !ph.proxyList.Add(p) {
		return
	}

	if ph.logs[p.Id] == nil {
		ph.logs[p.Id] = NewLogsList()
	}

	go func() {
		for {
			p.Run(ph.clientFilePath, func(time int64, _type, message string) {
				ph.logs[p.Id].Add(&Log{
					Time:    time,
					Type:    _type,
					Message: message,
				})
				ph.logsFunc(p.Id, time, _type, message)
			})
			var status cmd.Status
			select {
			case <-time.After(RUN_CONFIRM_DELAY):
				p.RunStatus = "running"
				ph.runStatusFunc(p.Id, "running")
				status = <-p.StatusChan
			case status = <-p.StatusChan:
			}
			log.Println("**********************************************************************************")
			log.Println(status)
			p.RunStatus = "idle"
			ph.runStatusFunc(p.Id, "idle")
			time.Sleep(DELAY_BEFORE_RUN_AGAIN)
			if p.Status == 0 {
				log.Println("ProxyHandler: AddProxy: breaking out of loop since p.Status is 0")
				break
			}
			if ph.ShutDown {
				log.Println("ProxyHandler: AddPrxoy: breaking out of loop since shutdown!")
				break
			}
		}
	}()
}

func (ph *ProxyHandler) RemoveProxy(id int64) {
	p := ph.proxyList.Remove(id)
	if p != nil {
		p.Status = 0
		p.Stop()
	}
}

func (ph *ProxyHandler) Start(ctx context.Context, dbHandler *DBHandler) {
	list, err := dbHandler.GetProxies(ctx, true)
	if err != nil {
		return
	}
	for i := range list {
		p := list[i]
		ph.AddProxy(&p)
	}
}

func (ph *ProxyHandler) Stop() {
	ph.ShutDown = true
	list := ph.proxyList.GetAll()
	for i := range list {
		p := list[i]
		if p != nil {
			p.Stop()
		}
	}
}

func (ph *ProxyHandler) GetLogs(id int64) []*Log {
	if ph.logs[id] == nil {
		return []*Log{}
	}
	return ph.logs[id].GetAll()
}

func (ph *ProxyHandler) DeleteLogs(id int64) {
	delete(ph.logs, id)
}

func NewProxyHandler(
	clientFilePath string,
	runStatusFunc ProxyRunStatusFunc,
	logsFunc LogsFunc,
) *ProxyHandler {
	return &ProxyHandler{
		proxyList:      NewProxyList(),
		clientFilePath: clientFilePath,
		runStatusFunc:  runStatusFunc,
		logsFunc:       logsFunc,
		logs:           map[int64]*LogsList{},
	}
}
