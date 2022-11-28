package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-cmd/cmd"
)

type Proxy struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	CreateTime int64  `json:"create_time"`
	S          string `json:"s"`
	P          int64  `json:"p"`
	B          string `json:"b"`
	L          int64  `json:"l"`
	K          string `json:"k"`
	M          string `json:"m"`
	O          string `json:"o"`
	Op         string `json:"op"`
	Oo         string `json:"oo"`
	Oop        string `json:"oop"`
	T          int64  `json:"t"`
	F          string `json:"f"`
	Status     int64  `json:"status"`

	RunStatus  string            `json:"run_status"`
	GoCmd      *cmd.Cmd          `json:"-"`
	StatusChan <-chan cmd.Status `json:"-"`
}

func (p *Proxy) GetCli() []string {
	return strings.Split(fmt.Sprintf(
		"-s %s -p %d -b %s -l %d -k %s -m %s -o %s --op %s -O %s --Op %s -t %d -f %s -v",
		p.S, p.P, p.B, p.L, p.K, p.M, p.O, p.Op, p.Oo, p.Oop, p.T, p.F), " ")
}

func (p *Proxy) Run(clientFilePath string, logsFun func(time int64, t string, l string)) {

	p.GoCmd = cmd.NewCmdOptions(cmd.Options{
		Buffered:  false,
		Streaming: true,
	}, clientFilePath, p.GetCli()...)
	p.StatusChan = p.GoCmd.Start()
	go func() {
		for p.GoCmd.Stdout != nil || p.GoCmd.Stderr != nil {
			select {
			case line, open := <-p.GoCmd.Stdout:
				if !open {
					p.GoCmd.Stdout = nil
					continue
				}
				logsFun(time.Now().UnixMilli(), "out", line)
			case line, open := <-p.GoCmd.Stderr:
				if !open {
					p.GoCmd.Stderr = nil
					continue
				}
				logsFun(time.Now().UnixMilli(), "err", line)
			}
		}
	}()
}

func (p *Proxy) Stop() {
	if p.GoCmd != nil {
		err := p.GoCmd.Stop()
		if err != nil {
			log.Println(err)
			err := KillProcessByPID(p.GoCmd.Status().PID)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
